package template

import (
	"os"
	"path/filepath"
)

func Apply(dir string, lang string) error {
	switch lang {
	case "python":
		return applyPython(dir)
	case "node":
		return applyNode(dir)
	case "go":
		return applyGo(dir)
	}
	return nil
}

func applyPython(dir string) error {
	name := filepath.Base(dir)
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("# "+name+"\n"), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "__init__.py"), []byte(""), 0644); err != nil {
		return err
	}
	main := `def main():
    print('Hello Python')

if __name__ == '__main__':
    main()
`
	if err := os.WriteFile(filepath.Join(dir, "main.py"), []byte(main), 0644); err != nil {
		return err
	}
	gitignore := `__pycache__/
*.py[cod]
.venv/
`
	return os.WriteFile(filepath.Join(dir, ".gitignore"), []byte(gitignore), 0644)
}

func applyNode(dir string) error {
	name := filepath.Base(dir)
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("# "+name+"\n"), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "index.js"), []byte("console.log('Hello Node');\n"), 0644); err != nil {
		return err
	}
	pkg := `{
  "name": "` + name + `",
  "version": "1.0.0",
  "main": "index.js",
  "license": "MIT"
}
`
	if err := os.WriteFile(filepath.Join(dir, "package.json"), []byte(pkg), 0644); err != nil {
		return err
	}
	gitignore := `node_modules/
.env
`
	return os.WriteFile(filepath.Join(dir, ".gitignore"), []byte(gitignore), 0644)
}

func applyGo(dir string) error {
	name := filepath.Base(dir)
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("# "+name+"\n"), 0644); err != nil {
		return err
	}
	main := `package main

import "fmt"

func main() {
	fmt.Println("Hello Go")
}
`
	return os.WriteFile(filepath.Join(dir, "main.go"), []byte(main), 0644)
}
