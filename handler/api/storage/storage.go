package storage

//func FindOneByUserId(storages core.StorageStore) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		user, ok := r.Context().Value(core.User{}).(*core.User)
//		if !ok {
//			w.WriteHeader(http.StatusInternalServerError)
//			log.Error("User are not inject in context")
//			return
//		}
//
//		storage, err := storages.GetByUserId(user.ID)
//
//		if err != nil {
//			w.WriteHeader(http.StatusInternalServerError)
//			log.Error(err)
//			return
//		} else if storage == nil {
//			http.NotFound(w, r)
//			return
//		}
//
//		out, err := json.Marshal(storage)
//		if err != nil {
//			w.WriteHeader(http.StatusInternalServerError)
//			log.Error(err)
//			return
//		}
//
//		w.WriteHeader(http.StatusOK)
//
//		_, err = w.Write(out)
//		if err != nil {
//			log.Error(err)
//			return
//		}
//	}
//}
