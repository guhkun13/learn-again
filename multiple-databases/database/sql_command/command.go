package sql_command

type Command struct {
	CreateTableIfNotExists string
}

func GetCommand() Command {
	return Command{
		CreateTableIfNotExists: "CREATE TABLE IF NOT EXISTS",
	}
}
