package request

type InitDB struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	UserName string `json:"userName" binding:"required"`
	Password string `json:"password"`
	DBName   string `json:"dbName" binding:"required"`
	// Admin login credentials (plain password; server stores MD5). Empty → gshark/gshark.
	AdminUserName string `json:"adminUserName"`
	AdminPassword string `json:"adminPassword"`
}
