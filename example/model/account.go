package model

type Account struct {
	ID            uint `gorm:"autoIncrement;primaryKey"`
	AccountName   string
	AccountNumber string
	Gender        string
	CellPhone     string
}

var AccountList = []Account{
	{
		AccountName:   "Bob",
		AccountNumber: "104938043525",
		Gender:        "male",
		CellPhone:     "+1 206-731-2132",
	},
	{
		AccountName:   "Billy",
		AccountNumber: "442897384221",
		Gender:        "male",
		CellPhone:     "+1 407-806-1432",
	},
	{
		AccountName:   "Mary",
		AccountNumber: "037733806596",
		Gender:        "female",
		CellPhone:     "+1 435-922-9518",
	},
	{
		AccountName:   "John",
		AccountNumber: "465179597739",
		Gender:        "male",
		CellPhone:     "+1 316-596-4621",
	},
	{
		AccountName:   "Jane",
		AccountNumber: "036073356959",
		Gender:        "female",
		CellPhone:     "+1 503-606-9341",
	},
	{
		AccountName:   "Alice",
		AccountNumber: "980055396477",
		Gender:        "female",
		CellPhone:     "+1 303-988-3603",
	},
	{
		AccountName:   "May",
		AccountNumber: "769369146010",
		Gender:        "female",
		CellPhone:     "+1 325-426-4328",
	},
	{
		AccountName:   "Austin",
		AccountNumber: "948322084003",
		Gender:        "male",
		CellPhone:     "+1 352-637-9490",
	},
	{
		AccountName:   "Jimmy",
		AccountNumber: "368834496660",
		Gender:        "male",
		CellPhone:     "+1 212-640-7230",
	},
	{
		AccountName:   "Jessica",
		AccountNumber: "154687370751",
		Gender:        "female",
		CellPhone:     "+1 213-812-5403",
	},
}
