package snowflake

var workerID *Snowflake

func init() {
	workerID, _ = NewSnowflake(0)
}

func Next() int64 {
	id, _ := workerID.Generate()
	return id
}
