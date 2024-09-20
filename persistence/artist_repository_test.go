package persistence

import (
	"context"

	"github.com/fatih/structs"
	"github.com/navidrome/navidrome/conf"
	"github.com/navidrome/navidrome/conf/configtest"
	"github.com/navidrome/navidrome/db"
	"github.com/navidrome/navidrome/log"
	"github.com/navidrome/navidrome/model"
	"github.com/navidrome/navidrome/model/request"
	"github.com/navidrome/navidrome/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("ArtistRepository", func() {
	var repo model.ArtistRepository

	BeforeEach(func() {
		ctx := log.NewContext(context.TODO())
		ctx = request.WithUser(ctx, model.User{ID: "userid"})
		repo = NewArtistRepository(ctx, NewDBXBuilder(db.Db()))
	})

	Describe("Count", func() {
		It("returns the number of artists in the DB", func() {
			Expect(repo.CountAll()).To(Equal(int64(2)))
		})
	})

	Describe("Exists", func() {
		It("returns true for an artist that is in the DB", func() {
			Expect(repo.Exists("3")).To(BeTrue())
		})
		It("returns false for an artist that is in the DB", func() {
			Expect(repo.Exists("666")).To(BeFalse())
		})
	})

	Describe("Get", func() {
		It("saves and retrieves data", func() {
			artist, err := repo.Get("2")
			Expect(err).ToNot(HaveOccurred())
			Expect(artist.Name).To(Equal(artistKraftwerk.Name))
		})
	})

	Describe("GetIndexKey", func() {
		// Note: OrderArtistName should never be empty, so we don't need to test for that
		r := artistRepository{indexGroups: utils.ParseIndexGroups(conf.Server.IndexGroups)}
		When("PreferSortTags is false", func() {
			BeforeEach(func() {
				DeferCleanup(configtest.SetupConfig)
				conf.Server.PreferSortTags = false
			})
			It("returns the OrderArtistName key is SortArtistName is empty", func() {
				conf.Server.PreferSortTags = false
				a := model.Artist{SortArtistName: "", OrderArtistName: "Bar", Name: "Qux"}
				idx := GetIndexKey(&r, a)
				Expect(idx).To(Equal("B"))
			})
			It("returns the OrderArtistName key even if SortArtistName is not empty", func() {
				a := model.Artist{SortArtistName: "Foo", OrderArtistName: "Bar", Name: "Qux"}
				idx := GetIndexKey(&r, a)
				Expect(idx).To(Equal("B"))
			})
		})
		When("PreferSortTags is true", func() {
			BeforeEach(func() {
				DeferCleanup(configtest.SetupConfig)
				conf.Server.PreferSortTags = true
			})
			It("returns the SortArtistName key if it is not empty", func() {
				a := model.Artist{SortArtistName: "Foo", OrderArtistName: "Bar", Name: "Qux"}
				idx := GetIndexKey(&r, a)
				Expect(idx).To(Equal("F"))
			})
			It("returns the OrderArtistName key if SortArtistName is empty", func() {
				a := model.Artist{SortArtistName: "", OrderArtistName: "Bar", Name: "Qux"}
				idx := GetIndexKey(&r, a)
				Expect(idx).To(Equal("B"))
			})
		})
	})

	Describe("GetIndex", func() {
		When("PreferSortTags is true", func() {
			BeforeEach(func() {
				DeferCleanup(configtest.SetupConfig())
				conf.Server.PreferSortTags = true
			})
			It("returns the index when PreferSortTags is true and SortArtistName is not empty", func() {
				// Set SortArtistName to "Foo" for Beatles
				artistBeatles.SortArtistName = "Foo"
				er := repo.Put(&artistBeatles)
				Expect(er).To(BeNil())

				idx, err := repo.GetIndex()
				Expect(err).ToNot(HaveOccurred())
				Expect(idx).To(HaveLen(2))
				Expect(idx[0].ID).To(Equal("F"))
				Expect(idx[0].Artists).To(HaveLen(1))
				Expect(idx[0].Artists[0].Name).To(Equal(artistBeatles.Name))
				Expect(idx[1].ID).To(Equal("K"))
				Expect(idx[1].Artists).To(HaveLen(1))
				Expect(idx[1].Artists[0].Name).To(Equal(artistKraftwerk.Name))

				// Restore the original value
				artistBeatles.SortArtistName = ""
				er = repo.Put(&artistBeatles)
				Expect(er).To(BeNil())
			})

			It("returns the index when PreferSortTags is true and SortArtistName is empty", func() {
				idx, err := repo.GetIndex()
				Expect(err).ToNot(HaveOccurred())
				Expect(idx).To(HaveLen(2))
				Expect(idx[0].ID).To(Equal("B"))
				Expect(idx[0].Artists).To(HaveLen(1))
				Expect(idx[0].Artists[0].Name).To(Equal(artistBeatles.Name))
				Expect(idx[1].ID).To(Equal("K"))
				Expect(idx[1].Artists).To(HaveLen(1))
				Expect(idx[1].Artists[0].Name).To(Equal(artistKraftwerk.Name))
			})
		})

		When("PreferSortTags is false", func() {
			BeforeEach(func() {
				DeferCleanup(configtest.SetupConfig())
				conf.Server.PreferSortTags = false
			})
			It("returns the index when SortArtistName is NOT empty", func() {
				// Set SortArtistName to "Foo" for Beatles
				artistBeatles.SortArtistName = "Foo"
				er := repo.Put(&artistBeatles)
				Expect(er).To(BeNil())

				idx, err := repo.GetIndex()
				Expect(err).ToNot(HaveOccurred())
				Expect(idx).To(HaveLen(2))
				Expect(idx[0].ID).To(Equal("B"))
				Expect(idx[0].Artists).To(HaveLen(1))
				Expect(idx[0].Artists[0].Name).To(Equal(artistBeatles.Name))
				Expect(idx[1].ID).To(Equal("K"))
				Expect(idx[1].Artists).To(HaveLen(1))
				Expect(idx[1].Artists[0].Name).To(Equal(artistKraftwerk.Name))

				// Restore the original value
				artistBeatles.SortArtistName = ""
				er = repo.Put(&artistBeatles)
				Expect(er).To(BeNil())
			})

			It("returns the index when SortArtistName is empty", func() {
				idx, err := repo.GetIndex()
				Expect(err).ToNot(HaveOccurred())
				Expect(idx).To(HaveLen(2))
				Expect(idx[0].ID).To(Equal("B"))
				Expect(idx[0].Artists).To(HaveLen(1))
				Expect(idx[0].Artists[0].Name).To(Equal(artistBeatles.Name))
				Expect(idx[1].ID).To(Equal("K"))
				Expect(idx[1].Artists).To(HaveLen(1))
				Expect(idx[1].Artists[0].Name).To(Equal(artistKraftwerk.Name))
			})
		})
	})

	Describe("dbArtist mapping", func() {
		var a *model.Artist
		BeforeEach(func() {
			a = &model.Artist{ID: "1", Name: "Van Halen", SimilarArtists: []model.Artist{
				{ID: "2", Name: "AC/DC"}, {ID: "-1", Name: "Test;With:Sep,Chars"},
			}}
		})
		It("maps fields", func() {
			dba := &dbArtist{Artist: a}
			m := structs.Map(dba)
			Expect(dba.PostMapArgs(m)).To(Succeed())
			Expect(m).To(HaveKeyWithValue("similar_artists", "2:AC%2FDC;-1:Test%3BWith%3ASep%2CChars"))

			other := dbArtist{SimilarArtists: m["similar_artists"].(string), Artist: &model.Artist{
				ID: "1", Name: "Van Halen",
			}}
			Expect(other.PostScan()).To(Succeed())

			actual := other.Artist
			Expect(*actual).To(MatchFields(IgnoreExtras, Fields{
				"ID":   Equal(a.ID),
				"Name": Equal(a.Name),
			}))
			Expect(actual.SimilarArtists).To(HaveLen(2))
			Expect(actual.SimilarArtists[0].ID).To(Equal("2"))
			Expect(actual.SimilarArtists[0].Name).To(Equal("AC/DC"))
			Expect(actual.SimilarArtists[1].ID).To(Equal("-1"))
			Expect(actual.SimilarArtists[1].Name).To(Equal("Test;With:Sep,Chars"))
		})
	})
})
