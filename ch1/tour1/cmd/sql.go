package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"tour1/internal/sql2struct"
)

var username string
var password string
var host string
var charset string
var dbType string
var dbName string
var tableName string

var sqlCmd = &cobra.Command{
	Use:   "sql",
	Short: "sql convert",
	Long:  "sql convert",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var sql2StructCmd = &cobra.Command{
	Use:   "struct",
	Short: "sql convert to struct",
	Long:  "sql convert to struct",
	Run: func(cmd *cobra.Command, args []string) {
		dbInfo := &sql2struct.DBInfo{
			DBType:   dbType,
			Host:     host,
			UserName: username,
			Password: password,
			Charset:  charset,
		}
		dbModel := sql2struct.NewDBModel(dbInfo)
		err := dbModel.Connect()

		if err != nil {
			log.Fatalf("dbModel connet err: %v\n", err)
		}
		columns, err := dbModel.GetColumns(dbName, tableName)
		if err != nil {
			log.Fatalf("dbModel.GetColumns err: %v\n", err)
		}
		template := sql2struct.NewStructTemplate()
		templateColumns := template.AssemblyColumns(columns)
		for i, c := range templateColumns {
			log.Printf("columns: %+v, %+v\n", columns[i], c)

		}

		err = template.Generate(tableName, templateColumns)
		if err != nil {
			log.Fatalf("template Generate err: %v\n", err)
		}
	},
}


func init() {
	sqlCmd.AddCommand(sql2StructCmd)
	sql2StructCmd.Flags().StringVarP(&username, "username", "u", "", "user name")
	sql2StructCmd.Flags().StringVarP(&password, "password", "p", "", "password")
	sql2StructCmd.Flags().StringVarP(&host, "host", "", "", "host")
	sql2StructCmd.Flags().StringVarP(&charset, "charset", "c", "", "char set")
	sql2StructCmd.Flags().StringVarP(&dbType, "type", "", "mysql", "database instance kind")
	sql2StructCmd.Flags().StringVarP(&dbName, "db", "", "", "database")
	sql2StructCmd.Flags().StringVarP(&tableName, "table", "", "", "table")
}