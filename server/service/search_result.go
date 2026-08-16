package service

import (
	"fmt"
	"strings"

	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/model"
	"github.com/madneal/gshark/model/request"
	"go.uber.org/zap"
)

// SaveResultStats contains detailed statistics about saved search results
type SaveResultStats struct {
	Total      int      // Total results processed
	Inserted   int      // Successfully inserted
	Skipped    int      // Skipped (already exists)
	Failed     int      // Failed to insert
	AIAnalyzed int      // Results sent to the AI pre-ingest filter
	AIRejected int      // Results classified as false positives
	AIErrors   int      // Results rejected because AI analysis failed
	Repos      []string // Unique repos affected
	repoSet    map[string]struct{}
}

// NewSaveResultStats creates a new SaveResultStats instance
func NewSaveResultStats() *SaveResultStats {
	return &SaveResultStats{
		Repos:   make([]string, 0),
		repoSet: make(map[string]struct{}),
	}
}

// AddRepo adds a repo to the stats if not already present
func (s *SaveResultStats) AddRepo(repo string) {
	if repo == "" {
		return
	}
	if _, exists := s.repoSet[repo]; !exists {
		s.repoSet[repo] = struct{}{}
		s.Repos = append(s.Repos, repo)
	}
}

// Summary returns a human-readable summary of the save operation
func (s *SaveResultStats) Summary(keyword, source string) string {
	aiSummary := ""
	if s.AIAnalyzed > 0 {
		aiSummary = fmt.Sprintf(", ai_analyzed=%d, ai_rejected=%d, ai_errors=%d", s.AIAnalyzed, s.AIRejected, s.AIErrors)
	}
	if s.Inserted == 0 {
		return fmt.Sprintf("[%s] keyword=%q: no new results (processed=%d, skipped=%d%s)",
			source, keyword, s.Total, s.Skipped, aiSummary)
	}

	repoSummary := ""
	if len(s.Repos) > 0 {
		if len(s.Repos) <= 3 {
			repoSummary = fmt.Sprintf(", repos=[%s]", strings.Join(s.Repos, ", "))
		} else {
			repoSummary = fmt.Sprintf(", repos=[%s, ...+%d more]",
				strings.Join(s.Repos[:3], ", "), len(s.Repos)-3)
		}
	}

	return fmt.Sprintf("[%s] keyword=%q: inserted=%d, skipped=%d, total=%d%s%s",
		source, keyword, s.Inserted, s.Skipped, s.Total, aiSummary, repoSummary)
}

func CreateSearchResult(searchResult model.SearchResult) (err error) {
	return Create(&searchResult)
}

func DeleteSearchResult(searchResult model.SearchResult) (err error) {
	return Delete(&searchResult)
}

func DeleteSearchResultByIds(ids request.IdsReq) (err error) {
	return DeleteByIds[model.SearchResult](ids)
}

func UpdateSearchResultByIds(req request.BatchUpdateReq) (err error) {
	err = global.GVA_DB.Table("search_result").Where("id in ?", req.Ids).
		UpdateColumn("status", req.Status).Error
	return err
}

func UpdateSearchResult(updateReq request.UpdateReq) (err error) {
	err = global.GVA_DB.Table("search_result").Where("repo = ?", updateReq.Repo).
		UpdateColumn("status", updateReq.Status).Error
	return err
}

func UpdateSearchResultById(id, status int) (err error) {
	err = global.GVA_DB.Table("search_result").Where("id = ?", id).
		UpdateColumn("status", status).Error
	return err
}

func GetSearchResult(id uint) (err error, searchResult model.SearchResult) {
	searchResult, err = GetByID[model.SearchResult](id)
	return
}

func ListSearchResultByStatus(status int) (err error, list []model.SearchResult) {
	err = global.GVA_DB.Where("status = ?", status).Find(&list).Error
	return err, list
}

func GetSearchResultInfoList(info request.SearchResultSearch) (err error, list interface{}, total int64) {
	db := global.GVA_DB.Model(&model.SearchResult{})
	var searchResults []model.SearchResult
	if info.Query != "" {
		db = db.Where("`repo` LIKE ? or `text_matches_json` LIKE ?",
			"%"+info.Query+"%", "%"+info.Query+"%")
	}
	if info.Keyword != "" {
		db = db.Where("`keyword` = ?", info.Keyword)
	}
	if info.Status >= 0 {
		db = db.Where("`status` = ?", info.Status)
	}
	total, err = Paginate(db, info.Page, info.PageSize, &searchResults, "id desc")
	return err, searchResults, total
}

func CheckExistOfSearchResult(searchResult *model.SearchResult) bool {
	urlExist := searchResult.CheckPathExists()
	repoExists := searchResult.CheckRepoExists()
	return urlExist || repoExists
}

func SaveSearchResults(searchResults []model.SearchResult) int {
	stats := SaveSearchResultsWithStats(searchResults)
	return stats.Inserted
}

func SaveSearchResultsWithStats(searchResults []model.SearchResult) *SaveResultStats {
	stats := NewSaveResultStats()
	stats.Total = len(searchResults)

	for _, result := range searchResults {
		exist := CheckExistOfSearchResult(&result)
		if exist {
			stats.Skipped++
			continue
		}
		if global.GVA_CONFIG.System.AiAnalysisEnabled {
			stats.AIAnalyzed++
			analysis, err := AnalyzeSearchResult(result)
			if err != nil {
				// AI filtering is deliberately fail-closed: when enabled, an
				// unavailable or malformed model response must not turn an
				// unverified finding into a stored finding.
				stats.AIErrors++
				global.GVA_LOG.Error("AI pre-ingest analysis failed; result rejected",
					zap.Uint("result_id", result.ID), zap.String("repo", result.Repo),
					zap.String("path", result.Path), zap.Error(err))
				continue
			}
			if !analysis.Real {
				stats.AIRejected++
				global.GVA_LOG.Info("AI pre-ingest analysis rejected result",
					zap.Uint("result_id", result.ID), zap.String("repo", result.Repo),
					zap.String("path", result.Path), zap.String("reason", analysis.Reason))
				continue
			}
		}
		err := CreateSearchResult(result)
		if err != nil {
			global.GVA_LOG.Error("save search result error", zap.Any("save searchResult error",
				err))
			stats.Failed++
		} else {
			stats.Inserted++
			stats.AddRepo(result.Repo)
		}
	}
	return stats
}

func SaveSearchResultPointers(searchResults []*model.SearchResult, keyword string) int {
	stats := SaveSearchResultPointersWithStats(searchResults, keyword)
	return stats.Inserted
}

func SaveSearchResultPointersWithStats(searchResults []*model.SearchResult, keyword string) *SaveResultStats {
	results := make([]model.SearchResult, 0, len(searchResults))
	for _, result := range searchResults {
		if result == nil {
			continue
		}
		if keyword != "" {
			result.Keyword = keyword
		}
		results = append(results, *result)
	}
	return SaveSearchResultsWithStats(results)
}
