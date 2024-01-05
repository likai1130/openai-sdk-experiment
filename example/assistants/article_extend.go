package assistants

type UserInfo struct {
	UserId            int    `json:"userId" gorm:"column:id;primaryKey;autoIncrement;comment:ID"`
	WalletAddr        string `json:"walletAddr" gorm:"not null;column:wallet_addr;comment:钱包地址"`
	Email             string `json:"email"  gorm:"not null;column:email;comment:邮箱"`
	NickName          string `json:"nickName" gorm:"column:nick_name;comment:昵称"`
	ProfileImage      string `json:"profileImage" gorm:"column:profile_image;comment:头像"`
	ProfileBackground string `json:"profileBackground" gorm:"column:profile_background;comment:头像背景色"`
	Bio               string `json:"bio" gorm:"type:text;column:bio;comment:个人经历"`
}

type PublicationDetails struct {
	Id          int    `json:"id" gorm:"column:id;primaryKey;autoIncrement;comment:ID"`
	Name        string `json:"name" gorm:"not null;;column:name;comment:期刊名称"`
	Description string `json:"description" gorm:"not null;column:description;comment:描述"`
	Logo        string `json:"logo" gorm:"column:logo;comment:图标"`
}

type DocumentInfo struct {
	Id            int    `json:"id" gorm:"column:id;primaryKey;autoIncrement;comment:ID"`
	PublicationId int    `json:"publicationId" gorm:"not null;default:0;column:publication_id;comment:期刊Id"`
	Content       string `json:"content" gorm:"column:content;comment:文章内容"`
	HtmlContent   string `json:"htmlContent" gorm:"column:html_content;comment:html类型的内容"`
	Version       string `json:"version" gorm:"column:version;comment:内容的版本"`
}

type TaskInfo struct {
	Id int `json:"id" gorm:"column:id;primaryKey;autoIncrement;comment:ID"`
}

// ArticleList 列表拼装展示对象
type ArticleList struct {
	Article
	UserInfo    UserInfo `json:"userInfo"`
	IsCollected bool     `json:"isCollected"` //是否收藏
}

// ArticleListDetails 拼装文章展示对象
type ArticleListDetails struct {
	ArticleList
	PublicationDetails PublicationDetails `json:"publicationDetails"`
}

// ArticleDetails 拼装文章展示对象
type ArticleDetails struct {
	ArticleListDetails
	DocumentInfo DocumentInfo `json:"documentInfo"`
	TaskInfo     TaskInfo     `json:"taskInfo"`
}

// BountyActiveInfo 分账日志信息
type BountyActiveInfo struct {
	Amount float64 `json:"amount" gorm:"column:amount;comment:分账金额"`
	TxHash string  `json:"tx_hash" gorm:"column:amount;comment:交易hash"`
}

// BountyWorkInfo 方案详情
type BountyWorkInfo struct {
	Article
	DocumentInfo DocumentInfo `json:"documentInfo" gorm:"comment:文档详情"`
}
