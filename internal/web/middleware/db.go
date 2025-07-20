package middleware

// const dbKey = "gormdb"

// // DB injects a *gorm.DB into the echo.Context
// func DB(db *db.DB) echo.MiddlewareFunc {
// 	return func(next echo.HandlerFunc) echo.HandlerFunc {
// 		return func(c echo.Context) error {
// 			c.Set(dbKey, db)
// 			return next(c)
// 		}
// 	}
// }

// // GetDB retrieves the *gorm.DB from echo.Context
// func GetDB(c echo.Context) *db.DB {
// 	if dbi, ok := c.Get(dbKey).(*db.DB); ok {
// 		return dbi
// 	}
// 	return nil
// }
