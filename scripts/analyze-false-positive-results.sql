-- Review-only report for simple token-prefix rules such as gho_ and ghp_.
-- Run against the GShark database; this script does not update any data.

SET @rule_keyword = 'gho_';
SET @minimum_ignored = 3;

-- Start with the distribution of manual decisions for the selected rule.
SELECT keyword, status, COUNT(*) AS result_count
FROM search_result
WHERE keyword = @rule_keyword
  AND status IN (1, 2)
GROUP BY keyword, status
ORDER BY status;

-- Extract repeated content signatures from text-match fragments. The token
-- value is normalized, while the surrounding line remains visible for review.
WITH fragments AS (
    SELECT
        result.status,
        matches.fragment
    FROM search_result AS result
    JOIN JSON_TABLE(
        CASE
            WHEN JSON_VALID(result.text_matches_json) THEN result.text_matches_json
            ELSE JSON_ARRAY()
        END,
        '$[*]' COLUMNS (
            fragment TEXT PATH '$.fragment'
        )
    ) AS matches ON TRUE
    WHERE result.keyword = @rule_keyword
      AND result.status IN (1, 2)
), normalized AS (
    SELECT
        status,
        LOWER(
            REGEXP_REPLACE(
                REGEXP_REPLACE(fragment, CONCAT(@rule_keyword, '[A-Za-z0-9_]*'), CONCAT(@rule_keyword, '<candidate>')),
                '[[:space:]]+',
                ' '
            )
        ) AS signature
    FROM fragments
    WHERE fragment IS NOT NULL
      AND fragment <> ''
      AND fragment LIKE CONCAT('%', @rule_keyword, '%')
)
SELECT
    signature,
    SUM(status = 2) AS ignored_count,
    SUM(status = 1) AS confirmed_count
FROM normalized
GROUP BY signature
HAVING ignored_count >= @minimum_ignored
   AND confirmed_count = 0
ORDER BY ignored_count DESC, signature;
