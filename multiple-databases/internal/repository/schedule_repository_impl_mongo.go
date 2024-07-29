package repository

// type ScheduleRepositoryImpl struct {
// 	Env    *config.EnvironmentVariable
// 	WrapDB *database.WrapDB
// }

// func NewScheduleRepositoryImpl(
// 	env *config.EnvironmentVariable,
// 	wrapDB *database.WrapDB,
// ) ScheduleRepository {
// 	return &ScheduleRepositoryImpl{
// 		Env:    env,
// 		WrapDB: wrapDB,
// 	}
// }

// func (r *ScheduleRepositoryImpl) Insert(req model.CronAutoCompress) (err error) {
// 	panic("implement this")

// 	// _, err = r.WrapDB.MongoApp.Collections.CronAutoCompress.InsertOne(context.TODO(), req)

// 	return
// }

// func (r *ScheduleRepositoryImpl) RemoveById(id string) (err error) {
// 	panic("implement this")

// 	// filter := bson.M{"id": id}
// 	// _, err = r.WrapDB.MongoApp.Collections.CronAutoCompress.DeleteOne(context.TODO(), filter)

// 	return
// }

// func (r *ScheduleRepositoryImpl) FindByCompressionId(id string) (res model.CronAutoCompress, err error) {
// 	panic("implement this")

// 	// filter := bson.M{"compression_id": id}
// 	// err = r.WrapDB.MongoApp.Collections.CronAutoCompress.FindOne(context.TODO(), filter).Decode(&res)

// 	return
// }

// func (r *ScheduleRepositoryImpl) UpdateCronEntryIdByCompressionId(compressionId string, entryId int) (err error) {
// 	panic("implement this")

// 	// isActive := false
// 	// if entryId > 0 {
// 	// 	isActive = true
// 	// }

// 	// filter := bson.M{"compression_id": compressionId}
// 	// update := bson.M{"$set": bson.M{"entry_id": entryId, "is_active": isActive}}

// 	// _, err = r.WrapDB.MongoApp.Collections.CronAutoCompress.UpdateOne(context.TODO(), filter, update)

// 	return
// }

// func (r *ScheduleRepositoryImpl) FindAll() (res model.CronAutoCompresses, err error) {
// 	panic("implement this")

// 	// filter := bson.M{}

// 	// cursor, err := r.WrapDB.MongoApp.Collections.CronAutoCompress.Find(context.TODO(), filter)
// 	// if err != nil {
// 	// 	return
// 	// }

// 	// err = cursor.All(context.TODO(), &res)

// 	return
// }
