package doc

var updateOperators = []string{
	"$currentData",
	"$inc",
	"$min",
	"$max",
	"$mul",
	"$rename",
	"$set",
	"$setOrInsert",
	"$unset",

	"$addToSet",
	"$pop",
	"$pull",
	"$push",
	"$pullAll",
}

type Update map[string]interface{}
