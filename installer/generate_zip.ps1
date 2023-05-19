Compress-Archive -Path ..\go\delete\Neo4jDelete -DestinationPath manual_install_files.zip -Force
Compress-Archive -Path ..\go\input\Neo4jInput -DestinationPath manual_install_files.zip -Update
Compress-Archive -Path ..\go\output\Neo4jOutput -DestinationPath manual_install_files.zip -Update
Compress-Archive -Path .\License.rtf -DestinationPath manual_install_files.zip -Update
Compress-Archive -Path .\graphyx.dll -DestinationPath manual_install_files.zip -Update
Compress-Archive -Path .\README.txt -DestinationPath manual_install_files.zip -Update