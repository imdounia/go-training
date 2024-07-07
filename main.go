package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	var choice string

	dbInfos := "root:@tcp(127.0.0.1:3306)/produits"

	db, err := sql.Open("mysql", dbInfos)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	for {
		fmt.Println("1. Ajouter un produit")
		fmt.Println("2. Afficher la liste des produits")
		fmt.Println("3. Modifier un produit")
		fmt.Println("4. Supprimer un produit")
		fmt.Println("5. Exporter les informations produits dans un fichier Excel (en .xlsx)")
		fmt.Println("6. Lancer un serveur Http avec une page web")
		fmt.Println("7. Se connecter à une VM en ssh")
		fmt.Println("8. Se connecter à un serveur FTP")
		fmt.Println("9. Lancer l'interface graphique")
		fmt.Println("10. Quitter")
		fmt.Println("Veuillez sélectionner une option :")

		fmt.Scan(&choice)

		switch choice {
		case "1":
			insertProduct(db)
		case "2":
			selectProducts(db)
		case "3":
			updateProduct(db)
		case "4":
			deleteProduct(db)
		case "5":
			exportProducts(db)
		case "6":
			return
		default:
			fmt.Println("Choix invalide, veuillez réessayer !")

		}
	}
}

func insertProduct(db *sql.DB) {
	var name, description, temp string
	var price float64

	fmt.Println("Veuillez entrer un nom :")
	fmt.Scan(&name)
	fmt.Println("Veuillez entrer une description :")
	fmt.Scan(&description)
	fmt.Println("Veuillez entrer un prix :")
	fmt.Scan(&temp)

	price, err := strconv.ParseFloat(temp, 64)
	if err != nil {
		fmt.Println("Erreur, veuillez renseigné un nombre flottant pour le prix !")
	}
	query := "INSERT INTO product (name, description, price) VALUES (?, ?, ?)"

	_, err = db.Exec(query, name, description, price)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Enregistrement réussi !")
}

func selectProducts(db *sql.DB) {
	query := "SELECT * FROM product"

	rows, err := db.Query(query)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		var id int
		var name, description string
		var price float64

		err = rows.Scan(&id, &name, &description, &price)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("ID : %d, Name : %s, Description : %s, Price : %f \n", id, name, description, price)

	}
}

func updateProduct(db *sql.DB) {
	var name, description, temp string
	var price float64

	fmt.Println("Renseignez l'Id du produit à modifier :")
	fmt.Scan(&temp)

	id, err := strconv.Atoi(temp)
	if err != nil {
		fmt.Println("Erreur : l'id doit être un entier ! ")
	}

	fmt.Println("Veuillez renseigner le nouveau nom :")
	fmt.Scan(&name)
	fmt.Println("Veuillez renseigner la nouvelle description :")
	fmt.Scan(&description)
	fmt.Println("Veuillez renseigner le nouveau prix :")
	fmt.Scan(&temp)

	price, err = strconv.ParseFloat(temp, 64)

	if err != nil {
		fmt.Println("Le prix doit être de type flottant !")
	}

	query := "UPDATE product SET name = ?, description = ?, price = ? WHERE id = ?"

	_, err = db.Exec(query, name, description, price, id)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Modification réussie !")
}

func deleteProduct(db *sql.DB) {
	var temp string
	var id int

	fmt.Println("Entrez l'id du produit à supprimer :")
	fmt.Scan(&temp)

	id, err := strconv.Atoi(temp)
	if err != nil {
		fmt.Println("L'id doit être un entier !")
	}

	query := "DELETE FROM product WHERE id = ?"

	_, err = db.Exec(query, id)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Suppression réussie !")
}

func exportProducts(db *sql.DB) {
	rows, err := db.Query("SELECT id, name, description, price FROM product")
	if err != nil {
		log.Fatalf("Erreur lors de la lecture des produits : %s", err)
	}
	defer rows.Close()

	xlsx := excelize.NewFile()
	sheetName := "Sheet1"
	xlsx.SetSheetName("Sheet1", sheetName)

	xlsx.SetCellValue(sheetName, "A1", "ID")
	xlsx.SetCellValue(sheetName, "B1", "Name")
	xlsx.SetCellValue(sheetName, "C1", "Description")
	xlsx.SetCellValue(sheetName, "D1", "Price")

	rowIndex := 2
	for rows.Next() {
		var id int
		var name, description string
		var price float64

		err := rows.Scan(&id, &name, &description, &price)
		if err != nil {
			log.Fatalf("Erreur lors du scan des données : %s", err)
		}

		xlsx.SetCellValue(sheetName, fmt.Sprintf("A%d", rowIndex), id)
		xlsx.SetCellValue(sheetName, fmt.Sprintf("B%d", rowIndex), name)
		xlsx.SetCellValue(sheetName, fmt.Sprintf("C%d", rowIndex), description)
		xlsx.SetCellValue(sheetName, fmt.Sprintf("D%d", rowIndex), price)

		rowIndex++
	}

	err = xlsx.SaveAs("products_export.xlsx")
	if err != nil {
		log.Fatalf("Erreur lors de la sauvegarde du fichier Excel : %s", err)
	}

	fmt.Println("Exportation des produits vers un fichier Excel réussie")
}
