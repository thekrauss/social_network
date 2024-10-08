package db

var (
	Users_table = `CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,  
		username TEXT UNIQUE NOT NULL,       
		email TEXT UNIQUE NOT NULL,          
		password_hash TEXT NOT NULL,          
		first_name TEXT NOT NULL,            
		last_name TEXT NOT NULL,              
		role TEXT CHECK( role IN ('admin', 'moderator', 'user') ) DEFAULT 'user',
		gender TEXT CHECK( gender IN ('male', 'female', 'other') ), 
		date_of_birth DATE NOT NULL,          
		avatar TEXT,                          
		bio TEXT,                             
		phone_number TEXT UNIQUE,             
		address TEXT,                         
		is_private BOOLEAN DEFAULT FALSE,     
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, 
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP  
	);`

	Followers_table = `CREATE TABLE IF NOT EXISTS followers (
		id TEXT PRIMARY KEY,
		follower_id TEXT NOT NULL,
		followed_id TEXT NOT NULL,
		status TEXT CHECK( status IN ('pending', 'accepted') ) DEFAULT 'pending',
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (follower_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (followed_id) REFERENCES users(id) ON DELETE CASCADE
	);`

	Groups_table = `CREATE TABLE IF NOT EXISTS groups (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT,
		creator_id TEXT NOT NULL,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (creator_id) REFERENCES users(id) ON DELETE CASCADE
	);`

	GroupMembers_table = `CREATE TABLE IF NOT EXISTS group_members (
		id TEXT PRIMARY KEY,
		group_id TEXT NOT NULL,
		user_id TEXT NOT NULL,
		status TEXT CHECK( status IN ('pending', 'accepted') ) DEFAULT 'pending',
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);`

	Notifications_table = `CREATE TABLE IF NOT EXISTS notifications (
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,          
		content TEXT NOT NULL,             
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		read BOOLEAN DEFAULT FALSE,        
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);`

	Messages_table = `CREATE TABLE IF NOT EXISTS messages (
		id TEXT PRIMARY KEY,
		sender_id TEXT NOT NULL,
		recipient_id TEXT NOT NULL,
		content TEXT NOT NULL,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (recipient_id) REFERENCES users(id) ON DELETE CASCADE
	);`

	Posts_table = `CREATE TABLE IF NOT EXISTS posts (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		user_id TEXT NOT NULL,
		visibility TEXT CHECK( visibility IN ('public', 'private', 'limited') ) DEFAULT 'public',
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		image_path TEXT,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);`

	Comments_table = `CREATE TABLE IF NOT EXISTS comments (
		id TEXT PRIMARY KEY,
		post_id TEXT NOT NULL,
		content TEXT NOT NULL,
		user_id TEXT NOT NULL,
		username TEXT NOT NULL,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);`

	CommentInteractions_table = `CREATE TABLE IF NOT EXISTS comment_interactions (
		id TEXT PRIMARY KEY,
		comment_id TEXT NOT NULL,
		user_id TEXT NOT NULL,
		interaction_type TEXT CHECK(interaction_type IN ('like', 'unlike')) NOT NULL,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE
	);`

	PostInteractions_table = `CREATE TABLE IF NOT EXISTS post_interactions (
		id TEXT PRIMARY KEY,
		post_id TEXT NOT NULL,
		user_id TEXT NOT NULL,
		interaction_type TEXT CHECK(interaction_type IN ('like', 'unlike')) NOT NULL,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
	);`
)
