package db

var (
	Users_tables = `CREATE TABLE IF NOT EXISTS users (
    	id TEXT PRIMARY KEY,  
		username TEXT UNIQUE NOT NULL,        -- Nom d'utilisateur unique
		email TEXT UNIQUE NOT NULL,           -- Adresse email unique
		password_hash TEXT NOT NULL,          -- Hachage du mot de passe (utiliser bcrypt ou une autre méthode de hachage)
		first_name TEXT NOT NULL,             -- Prénom de l'utilisateur
		last_name TEXT NOT NULL,              -- Nom de l'utilisateur
		role TEXT CHECK( role IN ('admin', 'moderator', 'user') ) DEFAULT 'user'
		gender TEXT CHECK( gender IN ('male', 'female', 'other') ), -- Genre de l'utilisateur (male, female, other)
		date_of_birth DATE NOT NULL,          -- Date de naissance
		avatar TEXT,                          -- URL ou chemin vers l'image d'avatar de l'utilisateur
		bio TEXT,                             -- Brève description ou biographie
		phone_number TEXT UNIQUE,             -- Numéro de téléphone
		address TEXT,                         -- Adresse de l'utilisateur (optionnel)
		is_private BOOLEAN DEFAULT FALSE,     -- Si le profil est privé ou non
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, -- Date de création du compte
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP  -- Dernière mise à jour du profil
	);`

	Followers_tables = `CREATE TABLE IF NOT EXISTS followers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		follower_id INTEGER NOT NULL,
		followed_id INTEGER NOT NULL,
		status TEXT CHECK( status IN ('pending', 'accepted') ) DEFAULT 'pending',
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (follower_id) REFERENCES users(id),
		FOREIGN KEY (followed_id) REFERENCES users(id)
	);`

	Groups_tables = `CREATE TABLE IF NOT EXISTS groups (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT,
		creator_id INTEGER NOT NULL,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (creator_id) REFERENCES users(id)
	);`

	Group_members_tables = `CREATE TABLE IF NOT EXISTS groupmembers (
    	id INTEGER PRIMARY KEY AUTOINCREMENT,
    	group_id INTEGER NOT NULL,
    	user_id INTEGER NOT NULL,
    	status TEXT CHECK( status IN ('pending', 'accepted') ) DEFAULT 'pending',
    	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    	FOREIGN KEY (group_id) REFERENCES groups(id),
    	FOREIGN KEY (user_id) REFERENCES users(id)
	);`

	Notifications = `CREATE TABLE IF NOT EXISTS notifications (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,          -- L'utilisateur qui reçoit la notification
		content TEXT NOT NULL,             -- Le contenu de la notification
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		read BOOLEAN DEFAULT FALSE,        -- Statut de lecture de la notification
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`

	Messages_table = `CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		sender_id INTEGER NOT NULL,
		recipient_id INTEGER NOT NULL,
		content TEXT NOT NULL,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (sender_id) REFERENCES users(id),
		FOREIGN KEY (recipient_id) REFERENCES users(id)
	);`

	Posts_table = `CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		category TEXT NOT NULL,
		content TEXT NOT NULL,
		user_id INTEGER NOT NULL,
		visibility TEXT CHECK( visibility IN ('public', 'private', 'limited') ) DEFAULT 'public',
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		image_path TEXT,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`

	CommentLikes_table = `CREATE TABLE IF NOT EXISTS commentlikes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		comment_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		FOREIGN KEY (comment_id) REFERENCES comments(id),
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`

	Comments_table = `CREATE TABLE IF NOT EXISTS comments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		post_id INTEGER NOT NULL,
		content TEXT NOT NULL,
		user_id INTEGER NOT NULL,
		username TEXT NOT NULL,
		created_at DATETIME NOT NULL,
		FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
	);`
)
