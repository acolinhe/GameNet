# GameNet

GameNet is a personal project aimed at creating a knowledge graph of interconnected video games.
The project utilizes Docker, Kubernetes, Go, PostgreSQL, Neo4j, and Natural Language Processing
(NLP) to fetch and process data from the Wikipedia API, storing the results in a structured
graph database.

## Introduction
GameNet aims to create a comprehensive knowledge graph that connects different video games
based on information fetched from Wikipedia. This project involves data ingestion, NLP
processing, and storing the processed data in a graph database for easy querying and analysis.

## Features
* Fetches data from the Wikipedia API
* Processes data using NLP techniques
* Stores data in PostgreSQL and Neo4j databases
* Utilizes Docker for containerization
* Deploys on Kubernetes for scalability

## Technologies
* **Go**: The main programming language for backend services
* **Docker**: For containerizing the application
* **Kubernetes**: For orchestrating the containers
* **PostgreSQL**: For relational data storage
* **Neo4j**: For graph data storage
* **NLP**: For processing and extracting information from text
* **Wikipedia API**: For fetching video game data

## Database

The GameNet project uses two databases, **PostgreSQL** and **Neo4j**, to manage video game articles and their associated metadata. Due to the sheer size of the dataset (thousands of video game articles and the relationships between them), it is impractical to store or host the database on GitHub. Below is an overview of the database structure and its contents.

### PostgreSQL

The **PostgreSQL** database serves as the primary relational database for storing structured data related to video games. It stores information such as:

- **Games**: Titles, summaries, and release dates.
- **Developers**: The companies or individuals who developed the games.
- **Genres**: The various genres each game falls under (e.g., action-adventure, platformer).
- **Platforms**: The gaming platforms (e.g., Nintendo Switch, PlayStation) the games are available on.

Each game is linked to multiple entities, such as developers, genres, and platforms. The relationships between these entities are stored in PostgreSQL using foreign keys, enabling efficient queries to retrieve metadata about the games.

With thousands of video game articles and their relationships, the database grows significantly. Each game can have multiple associated developers, platforms, and genres, resulting in a large amount of interrelated data. As a result, the database exceeds GitHub's size limits and cannot be stored directly in this repository.

### Neo4j

In addition to PostgreSQL, **Neo4j**, a graph database, is used to model the relationships between video games in a more flexible, graph-based format. This allows for querying relationships like:

- **Connections between developers and the games they have created**.
- **Games that share similar genres**.
- **Games that run on the same platforms**.

Neo4j is particularly useful for traversing relationships and discovering hidden patterns, such as finding common developers between different games or exploring games that belong to the same genre.

### Why Isn't the Database Stored on GitHub

The combined size of the PostgreSQL and Neo4j databases is too large to fit within GitHub’s repository limits. With thousands of video game articles, metadata, and relationships, the database requires external storage.

### How to Use the Database

To use the database locally:

1. **PostgreSQL**: The SQL schema for setting up the PostgreSQL database is included in this repository (`init.sql`). After setting up a PostgreSQL server, you can run the schema to create the necessary tables and relationships.

2. **Neo4j**: For Neo4j, the Cypher queries to populate the graph database are provided. You can set up a local or remote Neo4j instance and use these queries to import the relationships between video games.

Both databases are used in together to efficiently manage and query video game articles, their metadata, and the complex relationships between them.
