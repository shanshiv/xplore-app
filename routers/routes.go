package routers

import (
	"xplore/controller"

	"github.com/gorilla/mux"
)

func InitializeRoutes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/dashboard", controller.GetDashboard).Methods("GET")
	router.HandleFunc("/api/ulasan", controller.GetAllRatings).Methods("GET")
	router.HandleFunc("/api/beranda", controller.GetBeranda).Methods("GET")
	router.HandleFunc("/api/peta", controller.GetPeta).Methods("GET")
	router.HandleFunc("/api/transaksi", controller.CreateTransaksi).Methods("POST")
	router.HandleFunc("/api/transaksi/{id}", controller.GetTransaksiByUserID).Methods("GET")
	router.HandleFunc("/api/transaksi/sudah-diulas/{id}", controller.GetTransaksiWithReview).Methods("GET")
	router.HandleFunc("/api/transaksi/belum-diulas/{id}", controller.GetTransaksiWithoutReview).Methods("GET")

	router.HandleFunc("/api/favorit/{id}", controller.GetAllFavorit).Methods("GET")
	router.HandleFunc("/api/favoritkuliner/{id}", controller.GetAllKulinerFavorit).Methods("GET")
	router.HandleFunc("/api/favoritpenginapan/{id}", controller.GetAllPenginapanFavorit).Methods("GET")
	router.HandleFunc("/api/favoritwisata/{id}", controller.GetAllWisataFavorit).Methods("GET")
	router.HandleFunc("/api/user/kuliner/favorites", controller.AddKulinerFavorite).Methods("POST")
	router.HandleFunc("/api/user/penginapan/favorites", controller.AddPenginapanFavorite).Methods("POST")
	router.HandleFunc("/api/user/wisata/favorites", controller.AddWisataFavorite).Methods("POST")
	router.HandleFunc("/api/user/{user_id}/kuliner/favorites/{kuliner_id}", controller.DeleteKulinerFavorite).Methods("DELETE")
	router.HandleFunc("/api/user/{user_id}/penginapan/favorites/{penginapan_id}", controller.DeletePenginapanFavorite).Methods("DELETE")
	router.HandleFunc("/api/user/{user_id}/wisata/favorites/{wisata_id}", controller.DeleteWisataFavorite).Methods("DELETE")

	// Routes untuk akun
	router.HandleFunc("/api/loginadmin", controller.LoginAdmin).Methods("POST")
	router.HandleFunc("/api/loginuser", controller.LoginUser).Methods("POST")
	router.HandleFunc("/api/updateadminakun/{id}", controller.UpdateAdminAkun).Methods("PUT")
	router.HandleFunc("/api/akun", controller.CreateAkun).Methods("POST")
	router.HandleFunc("/api/akun", controller.GetAllAkun).Methods("GET")
	router.HandleFunc("/api/akun/{id}", controller.GetAkunByID).Methods("GET")
	router.HandleFunc("/api/akun/{id}", controller.UpdateAkun).Methods("PUT")
	router.HandleFunc("/api/akun/{id}", controller.DeleteAkun).Methods("DELETE")
	router.HandleFunc("/api/akun/email/{email}", controller.GetAkunByEmail).Methods("GET")

	// Routes untuk wisata
	router.HandleFunc("/api/wisata", controller.CreateWisata).Methods("POST")
	router.HandleFunc("/api/wisata", controller.GetAllWisata).Methods("GET")
	router.HandleFunc("/api/wisata/{id}", controller.GetWisataByID).Methods("GET")
	router.HandleFunc("/api/wisatadetail/{id}", controller.GetWisataDetailByID).Methods("GET")
	router.HandleFunc("/api/wisatadetailrating/{id}", controller.GetWisataRatingDetailByID).Methods("GET")
	router.HandleFunc("/api/wisata/{id}", controller.UpdateWisata).Methods("PUT")
	router.HandleFunc("/api/wisata/{id}", controller.DeleteWisata).Methods("DELETE")
	router.HandleFunc("/api/wisata_search", controller.SearchWisata).Methods("GET")
	router.HandleFunc("/api/wisata_searchrating", controller.SearchWisataRating).Methods("GET")
	router.HandleFunc("/api/wisatarating", controller.GetAllWisataRating).Methods("GET")
	//router.HandleFunc("/api/wisatapage", controller.GetWisataWithPagination).Methods("GET")
	// Routes untuk penginapan
	router.HandleFunc("/api/penginapan", controller.CreatePenginapan).Methods("POST")
	router.HandleFunc("/api/penginapan", controller.GetAllPenginapan).Methods("GET")
	router.HandleFunc("/api/penginapan/{id}", controller.GetPenginapanByID).Methods("GET")
	router.HandleFunc("/api/penginapandetail/{id}", controller.GetPenginapanDetailByID).Methods("GET")
	router.HandleFunc("/api/penginapandetailrating/{id}", controller.GetPenginapanRatingDetailByID).Methods("GET")
	router.HandleFunc("/api/penginapan/{id}", controller.UpdatePenginapan).Methods("PUT")
	router.HandleFunc("/api/penginapan/{id}", controller.DeletePenginapan).Methods("DELETE")
	router.HandleFunc("/api/penginapan_search", controller.SearchPenginapan).Methods("GET")
	router.HandleFunc("/api/penginapan_searchrating", controller.SearchPenginapanRating).Methods("GET")
	router.HandleFunc("/api/penginapanrating", controller.GetAllPenginapanRating).Methods("GET")
	// Routes untuk kuliner
	router.HandleFunc("/api/kuliner", controller.CreateKuliner).Methods("POST")
	router.HandleFunc("/api/kuliner", controller.GetAllKuliner).Methods("GET")
	router.HandleFunc("/api/kuliner/{id}", controller.GetKulinerByID).Methods("GET")
	router.HandleFunc("/api/kulinerdetail/{id}", controller.GetKulinerDetailByID).Methods("GET")
	router.HandleFunc("/api/kulinerdetailrating/{id}", controller.GetKulinerRatingDetailByID).Methods("GET")
	router.HandleFunc("/api/kuliner/{id}", controller.UpdateKuliner).Methods("PUT")
	router.HandleFunc("/api/kuliner/{id}", controller.DeleteKuliner).Methods("DELETE")
	router.HandleFunc("/api/kuliner_search", controller.SearchKuliner).Methods("GET")
	router.HandleFunc("/api/kuliner_searchrating", controller.SearchKulinerRating).Methods("GET")
	router.HandleFunc("/api/kulinerrating", controller.GetAllKulinerRating).Methods("GET")
	// Rute untuk rating wisata
	router.HandleFunc("/api/rating/wisata", controller.CreateRatingWisata).Methods("POST")
	router.HandleFunc("/api/rating/wisata", controller.GetAllRatingWisata).Methods("GET")
	router.HandleFunc("/api/rating/wisata/{id}", controller.GetRatingWisataByID).Methods("GET")
	router.HandleFunc("/api/rating/wisata/{id}", controller.UpdateRatingWisata).Methods("PUT")
	router.HandleFunc("/api/rating/wisata/{id}", controller.DeleteRatingWisata).Methods("DELETE")

	// Rute untuk rating kuliner
	router.HandleFunc("/api/rating/kuliner", controller.CreateRatingKuliner).Methods("POST")
	router.HandleFunc("/api/rating/kuliner", controller.GetAllRatingKuliner).Methods("GET")
	router.HandleFunc("/api/rating/kuliner/{id}", controller.GetRatingKulinerByID).Methods("GET")
	router.HandleFunc("/api/rating/kuliner/{id}", controller.UpdateRatingKuliner).Methods("PUT")
	router.HandleFunc("/api/rating/kuliner/{id}", controller.DeleteRatingKuliner).Methods("DELETE")

	// Rute untuk rating penginapan
	router.HandleFunc("/api/rating/penginapan", controller.CreateRatingPenginapan).Methods("POST")
	router.HandleFunc("/api/rating/penginapan", controller.GetAllRatingPenginapan).Methods("GET")
	router.HandleFunc("/api/rating/penginapan/{id}", controller.GetRatingPenginapanByID).Methods("GET")
	router.HandleFunc("/api/rating/penginapan/{id}", controller.UpdateRatingPenginapan).Methods("PUT")
	router.HandleFunc("/api/rating/penginapan/{id}", controller.DeleteRatingPenginapan).Methods("DELETE")

	return router
}
