package tests

// import (
//     "database/sql"
//     "testing"
// 	"damir/internal/entity"
// )


// func setupDB(t *testing.T) *sql.DB {
//     dsn := "postgres://postgres:12345@localhost:5432/fake?sslmode=disable"
//     db, err := sql.Open("postgres", dsn)
//     if err != nil {
//         t.Fatal("Failed to open a DB connection: ", err)
//     }
//     return db
// }

// func TestUserModelMethods(t *testing.T) {
//     db := setupDB(t)
//     defer db.Close()

//     _, err := db.Exec("DELETE FROM user_info")
//     if err != nil {
//         t.Fatalf("Failed to clear test database: %s", err)
//     }
	

//     userModel := &entity.UserModel{DB: db}
//     password := "securepassword123"
//     user := &entity.User{
//         Name:      "John",
//         Surname:   "Doe",
//         Email:     "john.doe@example.com",
//         Activated: true,
//     }

//     err = user.Password.Set(password)
//     if err != nil {
//         t.Fatalf("Failed to set password: %s", err)
//     }

//     err = userModel.Insert(user)
//     if err != nil {
//         t.Errorf("Failed to insert user: %s", err)
//     }

//     retrievedUser, err := userModel.Get(user.ID)
//     if err != nil {
//         t.Errorf("Failed to get user: %s", err)
//     } else if retrievedUser.Email != user.Email {
//         t.Errorf("Get user email %s does not match expected email %s", retrievedUser.Email, user.Email)
//     }

//     user.Name = "Updated John"
//     err = userModel.Update(user)
//     if err != nil {
//         t.Errorf("Failed to update user: %s", err)
//     }

//     updatedUser, err := userModel.Get(user.ID)
//     if err != nil {
//         t.Errorf("Failed to get user after update: %s", err)
//     } else if updatedUser.Name != "Updated John" {
//         t.Errorf("Update user name %s does not match expected name 'Updated John'", updatedUser.Name)
//     }

//     err = userModel.Delete(user.ID)
//     if err != nil {
//         t.Errorf("Failed to delete user: %s", err)
//     }
// }

// // func TestMovieModelMethods(t *testing.T) {
// //     db := setupDB(t)
// //     defer db.Close()

// //     _, err := db.Exec("DELETE FROM movies")
// //     if err != nil {
// //         t.Fatalf("Failed to clear test database: %s", err)
// //     }

// //     movieModel := &data.MovieModel{DB: db}
// //     movie := &data.Movie{
// //         Title:   "Inception",
// //         Year:    2010,
// //         Runtime: 148,
// //         Genres:  []string{"Action", "Sci-Fi", "Thriller"},
// //     }

// //     err = movieModel.Insert(movie)
// //     if err != nil {
// //         t.Fatalf("Failed to insert movie: %s", err)
// //     }

// //     retrievedMovie, err := movieModel.Get(movie.ID)
// //     if err != nil {
// //         t.Fatalf("Failed to get movie: %s", err)
// //     }
// //     if retrievedMovie.Title != movie.Title {
// //         t.Errorf("Retrieved movie title %s does not match expected title %s", retrievedMovie.Title, movie.Title)
// //     }

// //     movie.Title = "Inception Updated"
// //     err = movieModel.Update(movie)
// //     if err != nil {
// //         t.Fatalf("Failed to update movie: %s", err)
// //     }

// //     updatedMovie, err := movieModel.Get(movie.ID)
// //     if err != nil {
// //         t.Fatalf("Failed to get movie after update: %s", err)
// //     }
// //     if updatedMovie.Title != "Inception Updated" {
// //         t.Errorf("Updated movie title %s does not match expected title 'Inception Updated'", updatedMovie.Title)
// //     }

// //     err = movieModel.Delete(movie.ID)
// //     if err != nil {
// //         t.Fatalf("Failed to delete movie: %s", err)
// //     }

// //     _, err = movieModel.Get(movie.ID)
// //     if err == nil || err != data.ErrRecordNotFound {
// //         t.Errorf("Expected ErrRecordNotFound after movie deletion, got %v", err)
// //     }
// // }
