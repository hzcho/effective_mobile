package group

import (
	"net/http"
	"song_lib/internal/domain/model"
	"song_lib/internal/domain/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Song struct {
	songUsecase usecase.Song
}

func NewSong(songUsecases usecase.Song) *Song {
	return &Song{
		songUsecase: songUsecases,
	}
}

// @Summary Get a list of songs
// @Tags songs
// @Description Get a list of songs by filter with pagination
// @ID get-songs
// @Produce json
// @Param page query int false "Page number" default(0)
// @Param per_page query int false "Number of songs per page" default(10)
// @Param group query string false "Group name"
// @Param song query string false "Song title"
// @Success 200 {array} []model.SongDetails "List of songs"
// @Failure 400 {string} string "Invalid request format"
// @Failure 500 {string} string "Server error"
// @Router /api/v1/songs/info [get]
func (s *Song) GetLib(c *gin.Context) {
	pageStr := c.Query("page")
	perPageStr := c.Query("per_page")

	page := 0
	perPage := 0

	if pageStr != "" {
		var err error
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, "invalid page parameter")
			return
		}
	}

	if perPageStr != "" {
		var err error
		perPage, err = strconv.Atoi(perPageStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, "invalid per_page parameter")
			return
		}
	}

	group := c.Query("group")
	song := c.Query("song")

	input := model.LibraryFilter{
		Page:    page,
		PerPage: perPage,
		Group:   group,
		Song:    song,
	}

	songs, err := s.songUsecase.GetLib(c, input)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, "something went wrong")
		return
	}

	c.JSON(http.StatusOK, songs)
}

// @Summary Get song verses
// @Tags songs
// @Description Get verses of a specific song by ID with pagination
// @ID get-song-verses
// @Produce json
// @Param id path int true "Song ID"
// @Param page query int false "Page number" default(0)
// @Param per_page query int false "Number of verses per page" default(10)
// @Success 200 {array} model.VersesResponse "List of song verses"
// @Failure 400 {string} string "Invalid request format"
// @Failure 404 {string} string "Song not found"
// @Failure 500 {string} string "Server error"
// @Router /api/v1/songs/{id}/verses [get]
func (s *Song) GetVerses(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid song ID")
		return
	}

	verseStr := c.Query("page")
	perVerseStr := c.Query("per_page")

	verse := 0
	perVerse := 0

	if verseStr != "" {
		var err error
		verse, err = strconv.Atoi(verseStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, "invalid page parameter")
			return
		}
	}

	if perVerseStr != "" {
		var err error
		perVerse, err = strconv.Atoi(perVerseStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, "invalid per_page parameter")
			return
		}
	}

	input := model.VersesRequest{
		SongID:  id,
		Page:    verse,
		PerPage: perVerse,
	}

	verses, err := s.songUsecase.GetVerses(c, input)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, "something went wrong")
		return
	}

	c.JSON(http.StatusOK, verses)
}

// @Summary Add a new song
// @Tags songs
// @Description Add a new song to the library
// @ID add-song
// @Produce json
// @Param song body model.AddSong true "Song details"
// @Success 200 {integer} int "ID of the created song"
// @Failure 400 {string} string "Incorrect fields"
// @Failure 500 {string} string "Server error"
// @Router /api/v1/songs [post]
func (s *Song) Add(c *gin.Context) {
	input := model.AddSong{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, "incorrect fields")
		return
	}

	id, err := s.songUsecase.Add(c, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "something went wrong")
		return
	}

	c.JSON(http.StatusOK, id)
}

// @Summary Update an existing song
// @Tags songs
// @Description Update the details of an existing song by ID
// @ID update-song
// @Produce json
// @Param id path int true "Song ID"
// @Param song body model.UpdateSong true "Updated song details"
// @Success 200 {object} model.Song "Updated song details"
// @Failure 400 {string} string "Incorrect fields"
// @Failure 404 {string} string "Song not found"
// @Failure 500 {string} string "Server error"
// @Router /api/v1/songs/{id} [put]
func (s *Song) Update(c *gin.Context) {
	input := model.UpdateSong{}

	if err := c.Bind(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "incorrect fields")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid song ID")
		return
	}

	input.ID = id

	song, err := s.songUsecase.Update(c, input)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, "something went wrong")
		return
	}

	c.JSON(http.StatusOK, song)
}

// @Summary Delete a song
// @Tags songs
// @Description Delete a song from the library by ID
// @ID delete-song
// @Produce json
// @Param id path int true "Song ID"
// @Success 200 {string} string "The song has been deleted"
// @Failure 400 {string} string "Invalid song ID or something went wrong"
// @Failure 404 {string} string "Song not found"
// @Failure 500 {string} string "Server error"
// @Router /api/v1/songs/{id} [delete]
func (s *Song) Delete(c *gin.Context) {

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid song ID")
		return
	}

	err = s.songUsecase.Delete(c, id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "something went wrong")
		return
	}

	c.JSON(http.StatusOK, "the song has been deleted")
}
