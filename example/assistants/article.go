package assistants

type Article struct {
	Id            int    `json:"id" gorm:"column:id;primaryKey;autoIncrement;comment:ID"`
	PublicationId int    `json:"publicationId" gorm:"not null;default:0;column:publication_id;comment:期刊Id"`
	BountyId      int    `json:"bountyId" gorm:"column:bounty_id;comment:征稿Id"`
	Type          string `json:"type" gorm:"default:'';column:type;comment:确定是论文还是征稿"`
	Title         string `json:"title" gorm:"not null;;column:title;comment:标题"`
	Abstract      string `json:"abstract" gorm:"column:abstract;comment:摘要信息"`
	Directory     string `json:"directory" gorm:"column:directory;comment:文章目录"`
	Status        int    `json:"status" gorm:"column:status;comment: 文章状态: 1:草稿 2:审核中 3:手动撤回 4：被退稿 5：通过"`
	ContentCid    string `json:"contentCid" gorm:"default:'';column:content_cid;comment:内容cid"`
	MetaCid       string `json:"metaCid" gorm:"column:meta_cid;comment:元数据cid"`
	DocId         int    `json:"docId" gorm:"column:doc_id;comment:文档ID"`
	UserId        int    `json:"userId" gorm:"column:user_id;comment:用户ID"`
	IsWinner      bool   `json:"isWinner" gorm:"column:is_winner;comment:获胜者"`
}
