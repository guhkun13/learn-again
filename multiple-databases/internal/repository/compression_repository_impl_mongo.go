package repository

// func (r *CompressionRepositoryImpl) MongoInsert(req model.Compression) error {
// 	_, err := r.WrapDB.MongoApp.Collections.Compression.InsertOne(context.TODO(), req)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (r *CompressionRepositoryImpl) MongoUpdateAutoCompressPath(compression model.Compression, req dto.UpdateAutoCompressRequest) error {
// 	filter := bson.M{"id": compression.Id}
// 	update := bson.M{"$set": bson.M{"source_path": req.SourcePath, "destination_path": req.DestinationPath}}

// 	_, err := r.WrapDB.MongoApp.Collections.Compression.UpdateOne(context.TODO(), filter, update)

// 	return err
// }

// func (r *CompressionRepositoryImpl) MongoFindAll(f dto.FilterFindCompression) (res model.FindAllCompressionResult, err error) {
// 	filter := bson.M{}

// 	if f.Repetition != "" {
// 		filter["repetition"] = f.Repetition
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), r.Env.Database.Timeout.Read)
// 	defer cancel()

// 	var compressions model.Compressions

// 	pgData, err := mongopagination.New(r.WrapDB.MongoApp.Collections.Compression).Context(ctx).
// 		Limit(int64(f.Pagination.Limit)).
// 		Page(int64(f.Pagination.Page)).
// 		Sort(f.Sort.SortBy, helper.ConvertSortOrderDirectionForMongo(f.Sort.SortOrder)).
// 		Filter(filter).
// 		Decode(&compressions).Find()

// 	if err != nil {
// 		log.Error().Err(err).Msg("CompressionRepositoryImpl.FindAll failed")
// 		return
// 	}

// 	res.Compressions = compressions
// 	res.PaginatedData = pgData

// 	return
// }

// func (r *CompressionRepositoryImpl) MongoFindById(id string) (res model.Compression, err error) {
// 	filter := bson.M{"id": id}
// 	err = r.WrapDB.MongoApp.Collections.Compression.FindOne(context.TODO(), filter).Decode(&res)
// 	if err != nil {
// 		return
// 	}

// 	return
// }

// func (r *CompressionRepositoryImpl) MongoRemoveById(id string) (err error) {
// 	filter := bson.M{"id": id}
// 	_, err = r.WrapDB.MongoApp.Collections.Compression.DeleteOne(context.TODO(), filter)
// 	if err != nil {
// 		return
// 	}

// 	return
// }

// // Multiple

// func (r *CompressionRepositoryImpl) MongoFindByIds(ids []string) (res model.Compressions, err error) {
// 	filter := bson.M{"id": bson.M{"$in": ids}}
// 	cursor, err := r.WrapDB.MongoApp.Collections.Compression.Find(context.TODO(), filter)
// 	if err != nil {
// 		return
// 	}

// 	err = cursor.All(context.TODO(), &res)

// 	return
// }

// func (r *CompressionRepositoryImpl) MongoDeleteByIds(ids []string) (countDeleted int, err error) {
// 	filter := bson.M{"id": bson.M{"$in": ids}}
// 	result, err := r.WrapDB.MongoApp.Collections.Compression.DeleteMany(context.TODO(), filter)
// 	if err != nil {
// 		return
// 	}

// 	log.Info().Int("records deleted", int(result.DeletedCount)).Msg("result")
// 	countDeleted = int(result.DeletedCount)

// 	return
// }
