import os
directory = os.getcwd()
output_file = os.path.join(directory, "output.txt")

# Lists of folder names and file names to ignore
folders_to_ignore = ['blockchain_data', "summaries", ".git"]
files_to_ignore = ['prompt.txt', 'output.txt', "read.py", "resp.txt","go.mod", "go.sum", ".gitignore", "package-lock.json", "README.md", "app.check", "p2p-server.check","go.mod","go.sum"]

with open(output_file, 'w', encoding='utf-8') as out_file:
    for root, dirs, files in os.walk(directory):
        # Remove ignored folders from the search
        dirs[:] = [d for d in dirs if d not in folders_to_ignore]
        
        for file in files:
            # Skip ignored files
            if file in files_to_ignore:
                continue
            
            file_path = os.path.join(root, file)
            out_file.write(f"File: {file_path}\n")
            try:
                with open(file_path, 'r', encoding='utf-8') as f:
                    out_file.write(f.read())
            except Exception as e:
                out_file.write(f"Error reading {file_path}: {e}\n")
            out_file.write("\n")

print(f"Done! Output saved to {output_file}")