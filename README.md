# graphyx

Graphyx is a set of Neo4j connectors for Alteryx:

* Neo4j Input: Import cypher queries into Alteryx workflows
* Neo4j Output: Export Alteryx data as Neo4j nodes and relationships
* Neo4j Delete: Use Alteryx data to define how Neo4j nodes and relationships should be deleted

The engine for the connectors was built using the [Alteryx Go SDK](https://github.com/tlarsen7572/goalteryx) and the [official Go driver for Neo4j](https://github.com/neo4j/neo4j-go-driver).

The user interfaces were built using [Flutter](https://github.com/flutter/flutter).

## Table of contents

1. [Installation](#Installation)
2. [Neo4j Input](#Neo4j-Input)
3. [Neo4j Output](#Neo4j-Output)
4. [Neo4j Delete](#Neo4j-Delete)

## Installation

1. Download graphyx.zip from the latest [release](https://github.com/tlarsen7572/graphyx/releases).
2. Extract the zip file to a temporary location.
3. Copy the Neo4jInput, Neo4jOutput, and Neo4jDelete folders to an appropriate custom tool folder. To install for all users, the path is typically `C:\ProgramData\Alteryx\Tools`. For a user-specific install, the path is typically `C:\Users\YourUsername\AppData\Roaming\Alteryx\Tools`. This step installs the tool configurations and UIs.
4. Copy graphyx.dll to your Alteryx plugins folder. For an admin installation of Alteryx, the path is typically `C:\Program Files\Alteryx\bin\Plugins`. For a user-specific installation, the path is typically `C:\Users\YourUsername\AppData\Local\Alteryx\bin\Plugins`. This step installs the engine for the connectors.

Graphyx is now installed. You can find the new connectors in the Connectors tab in Designer.

[Back to top](#graphyx)

## Neo4j Input

<img src="https://github.com/tlarsen7572/graphyx/blob/main/go/input/Neo4jInput/icon.png" width="100" />

[Back to top](#graphyx)

## Neo4j Output

<img src="https://github.com/tlarsen7572/graphyx/blob/main/go/output/Neo4jOutput/icon.png" width="100" />

[Back to top](#graphyx)

## Neo4j Delete

<img src="https://github.com/tlarsen7572/graphyx/blob/main/go/delete/Neo4jDelete/icon.png" width="100" />

[Back to top](#graphyx)
