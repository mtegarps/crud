package main

import (
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "net/http"
    "time"
	"strconv"
    "strings"
)

type Contact struct {
    ID        string    `json:"id"`
    Name      string    `json:"name"`
    Gender    string    `json:"gender"`
    Phone     string    `json:"phone"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

var contacts []Contact

func main() {
    r := gin.Default()

    // Create a contact
    r.POST("/contacts", createContact)

    // Retrieve all contacts
    r.GET("/contacts", getAllContacts)

    // Retrieve a contact by ID
    r.GET("/contacts/:id", getContactByID)

    // Update a contact by ID
    r.PUT("/contacts/:id", updateContact)

    // Delete a contact by ID
    r.DELETE("/contacts/:id", deleteContact)

    r.Run(":8080")
}

func createContact(c *gin.Context) {
    var contact Contact
    if err := c.ShouldBindJSON(&contact); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    contact.ID = uuid.New().String()
    contact.CreatedAt = time.Now()
    contact.UpdatedAt = time.Now()

    contacts = append(contacts, contact)

    c.JSON(http.StatusCreated, contact)
}

func getAllContacts(c *gin.Context) {
    // Filter by name (query parameter)
    nameFilter := c.Query("name")

    // Pagination parameters (query parameters)
    page := c.DefaultQuery("page", "1")
    pageSize := c.DefaultQuery("page_size", "10")

    // Convert page and pageSize to integers
    pageInt, _ := strconv.Atoi(page)
    pageSizeInt, _ := strconv.Atoi(pageSize)

    // Apply filter if nameFilter is provided
    filteredContacts := contacts
    if nameFilter != "" {
        filteredContacts = filterContactsByName(filteredContacts, nameFilter)
    }

    // Calculate total count and paginate the filtered contacts
    totalCount := len(filteredContacts)
    paginatedContacts := paginateContacts(filteredContacts, pageInt, pageSizeInt)

    // Return the paginated contacts as JSON response
    c.JSON(http.StatusOK, gin.H{
        "total_count": totalCount,
        "page":        pageInt,
        "page_size":   pageSizeInt,
        "contacts":    paginatedContacts,
    })
}

func filterContactsByName(contacts []Contact, nameFilter string) []Contact {
    filtered := []Contact{}
    for _, contact := range contacts {
        if strings.Contains(contact.Name, nameFilter) {
            filtered = append(filtered, contact)
        }
    }
    return filtered
}

func paginateContacts(contacts []Contact, page, pageSize int) []Contact {
    startIndex := (page - 1) * pageSize
    if startIndex >= len(contacts) {
        return []Contact{}
    }

    endIndex := startIndex + pageSize
    if endIndex > len(contacts) {
        endIndex = len(contacts)
    }

    return contacts[startIndex:endIndex]
}

func getContactByID(c *gin.Context) {
    id := c.Param("id")

    for _, contact := range contacts {
        if contact.ID == id {
            c.JSON(http.StatusOK, contact)
            return
        }
    }

    c.JSON(http.StatusNotFound, gin.H{"error": "Contact not found"})
}

func updateContact(c *gin.Context) {
    id := c.Param("id")

    var updatedContact Contact
    if err := c.ShouldBindJSON(&updatedContact); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var foundContact *Contact
    for i, contact := range contacts {
        if contact.ID == id {
            foundContact = &contacts[i]
            break
        }
    }

    if foundContact == nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Contact not found"})
        return
    }

    foundContact.Name = updatedContact.Name
    foundContact.Gender = updatedContact.Gender
    foundContact.Phone = updatedContact.Phone
    foundContact.Email = updatedContact.Email
    foundContact.UpdatedAt = time.Now()

    c.JSON(http.StatusOK, foundContact)
}

func deleteContact(c *gin.Context) {
    id := c.Param("id")

    for i, contact := range contacts {
        if contact.ID == id {
            contacts = append(contacts[:i], contacts[i+1:]...)
            c.JSON(http.StatusOK, gin.H{"message": "Contact deleted"})
            return
        }
    }

    c.JSON(http.StatusNotFound, gin.H{"error": "Contact not found"})
}
