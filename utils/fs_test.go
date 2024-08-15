package utils

import (
	"os"
	"testing"
)

func TestPathExists(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	tests := []struct {
		path   string
		exists bool
	}{
		{tmpFile.Name(), true},
		{"nonexistentfile.txt", false},
	}

	for _, test := range tests {
		result := PathExists(test.path)
		if result != test.exists {
			t.Errorf("PathExists(%q) = %v; expected %v", test.path, result, test.exists)
		}
	}
}

func TestMakeDir(t *testing.T) {
	testDir := "./testdir/subdir"

	defer func() {
		if err := os.RemoveAll("./testdir"); err != nil {
			t.Fatalf("Failed to clean up test directory: %v", err)
		}
	}()

	err := MakeDir(testDir)
	if err != nil {
		t.Errorf("MakeDir(%q) returned an error: %v", testDir, err)
		return
	}

	if _, err := os.Stat(testDir); os.IsNotExist(err) {
		t.Errorf("Directory %q was not created", testDir)
	}
}

func TestWriteToFile(t *testing.T) {
	testFile := "./testfile.txt"
	defer func() {
		if err := os.Remove(testFile); err != nil {
			t.Fatalf("Failed to clean up test file: %v", err)
		}
	}()

	tests := []struct {
		content   string
		expectErr bool
	}{
		{"Hello, world!", false}, // Write new content
		{"Hello, world!", false},
		{"New content!", false},
		{"", false},
	}

	for _, test := range tests {
		if err := os.WriteFile(testFile, []byte("Initial content"), 0644); err != nil {
			t.Fatalf("Failed to set up initial file content: %v", err)
		}
		err := WriteToFile(testFile, test.content)
		if (err != nil) != test.expectErr {
			t.Errorf("WriteToFile(%q) returned error: %v; expected error: %v", test.content, err, test.expectErr)
			continue
		}

		if err == nil {
			data, err := os.ReadFile(testFile)
			if err != nil {
				t.Errorf("Failed to read file content: %v", err)
				continue
			}

			if string(data) != test.content {
				t.Errorf("File content = %q; expected %q", string(data), test.content)
			}
		}
	}
}

func TestRemovePath(t *testing.T) {
	testDir := "./testdir"
	testFile := "./testdir/testfile.txt"
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Fatalf("Test file does not exist before RemovePath")
	}
	if _, err := os.Stat(testDir); os.IsNotExist(err) {
		t.Fatalf("Test directory does not exist before RemovePath")
	}
	err := RemovePath(testDir)
	if err != nil {
		t.Errorf("RemovePath(%q) returned an error: %v", testDir, err)
		return
	}
	if _, err := os.Stat(testFile); !os.IsNotExist(err) {
		t.Errorf("Test file still exists after RemovePath")
	}
	if _, err := os.Stat(testDir); !os.IsNotExist(err) {
		t.Errorf("Test directory still exists after RemovePath")
	}
}

func TestMovePath(t *testing.T) {
	oldDir := "./oldpath"
	oldPath := oldDir + "/testfile.txt"
	newDir := "./newpath"
	newPath := newDir + "/testfile.txt"

	if err := os.MkdirAll(oldDir, 0755); err != nil {
		t.Fatalf("Failed to create oldpath directory: %v", err)
	}
	if err := os.WriteFile(oldPath, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	if _, err := os.Stat(oldPath); os.IsNotExist(err) {
		t.Fatalf("Test file does not exist at old path before MovePath")
	}

	if err := os.MkdirAll(newDir, 0755); err != nil {
		t.Fatalf("Failed to create newpath directory: %v", err)
	}

	err := MovePath(oldPath, newPath)
	if err != nil {
		t.Errorf("MovePath(%q, %q) returned an error: %v", oldPath, newPath, err)
		return
	}

	if _, err := os.Stat(newPath); os.IsNotExist(err) {
		t.Errorf("Test file does not exist at new path after MovePath")
	}
	if _, err := os.Stat(oldPath); !os.IsNotExist(err) {
		t.Errorf("Test file still exists at old path after MovePath")
	}

	if err := os.Remove(newPath); err != nil {
		t.Errorf("Failed to clean up test file at new path: %v", err)
	}
	if err := os.RemoveAll(newDir); err != nil {
		t.Errorf("Failed to clean up newpath directory: %v", err)
	}
	if err := os.RemoveAll(oldDir); err != nil {
		t.Errorf("Failed to clean up oldpath directory: %v", err)
	}
}

func TestCopyFile(t *testing.T) {
	srcFile := "./sourcefile.txt"
	dstFile := "./destinationfile.txt"
	if err := os.WriteFile(srcFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}
	if _, err := os.Stat(srcFile); os.IsNotExist(err) {
		t.Fatalf("Source file does not exist before CopyFile")
	}
	err := CopyFile(srcFile, dstFile)
	if err != nil {
		t.Errorf("CopyFile(%q, %q) returned an error: %v", srcFile, dstFile, err)
		return
	}
	destData, err := os.ReadFile(dstFile)
	if err != nil {
		t.Errorf("Failed to read destination file: %v", err)
		return
	}

	srcData, err := os.ReadFile(srcFile)
	if err != nil {
		t.Errorf("Failed to read source file: %v", err)
		return
	}

	if string(destData) != string(srcData) {
		t.Errorf("Destination file content = %q; expected %q", string(destData), string(srcData))
	}
	if err := os.Remove(srcFile); err != nil {
		t.Errorf("Failed to clean up source file: %v", err)
	}
	if err := os.Remove(dstFile); err != nil {
		t.Errorf("Failed to clean up destination file: %v", err)
	}
}

func TestReplaceInFile(t *testing.T) {
	filePath := "./testfile.txt"
	initialContent := "Hello, world! Hello again!"
	if err := os.WriteFile(filePath, []byte(initialContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	oldStr := "Hello"
	newStr := "Hi"
	err := ReplaceInFile(filePath, oldStr, newStr)
	if err != nil {
		t.Errorf("ReplaceInFile(%q, %q, %q) returned an error: %v", filePath, oldStr, newStr, err)
		return
	}
	expectedContent := "Hi, world! Hi again!"
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Errorf("Failed to read file content: %v", err)
		return
	}

	if string(data) != expectedContent {
		t.Errorf("File content = %q; expected %q", string(data), expectedContent)
	}
	if err := os.Remove(filePath); err != nil {
		t.Errorf("Failed to clean up test file: %v", err)
	}
}

func TestReplaceManyInFile(t *testing.T) {
	filePath := "./testfile.txt"
	initialContent := "Hello, world! Goodbye, world!"
	if err := os.WriteFile(filePath, []byte(initialContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	replacements := map[string]string{
		"Hello":   "Hi",
		"Goodbye": "Farewell",
	}
	err := ReplaceManyInFile(filePath, replacements)
	if err != nil {
		t.Errorf("ReplaceManyInFile(%q, %v) returned an error: %v", filePath, replacements, err)
		return
	}
	expectedContent := "Hi, world! Farewell, world!"
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Errorf("Failed to read file content: %v", err)
		return
	}

	if string(data) != expectedContent {
		t.Errorf("File content = %q; expected %q", string(data), expectedContent)
	}
	if err := os.Remove(filePath); err != nil {
		t.Errorf("Failed to clean up test file: %v", err)
	}
}
