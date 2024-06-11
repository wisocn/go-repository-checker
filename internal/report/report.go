package report

import (
	"fmt"
	"os"
	"text/template"
)

func generateReport(validations []Validation) {
	fmt.Println("Generating report")

	// Read the HTML template from an external file
	templateFile := "report-template.html"
	tmplContent, err := os.ReadFile(templateFile)
	if err != nil {
		fmt.Println("Error reading template file:", err)
		return
	}

	// Create an HTML file
	file, err := os.Create("report.html")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Parse and execute the template
	tmpl, err := template.New("report").Parse(string(tmplContent))
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return
	}

	err = tmpl.Execute(file, validations)
	if err != nil {
		fmt.Println("Error executing template:", err)
		return
	}

	fmt.Println("Report generated successfully.")

}
