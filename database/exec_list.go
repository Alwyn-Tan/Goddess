package database

func execLPush(database *Database, args [][]byte) redis.Reply {
	key := string(args[0])
	values := args[1:]

	list, _, errReply := database.getOrInitList(key)

}
