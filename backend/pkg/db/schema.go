package db

var (
	Users_tables = `CREATE TABLE IF NOT EXISTS users (
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

	Followers_tables = `CREATE TABLE IF NOT EXISTS followers (
			id TEXT PRIMARY KEY,
			follower_id TEXT NOT NULL,
			followed_id TEXT NOT NULL,
			status TEXT CHECK( status IN ('pending', 'accepted') ) DEFAULT 'pending',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (follower_id) REFERENCES users(id),
			FOREIGN KEY (followed_id) REFERENCES users(id)
		);`

	Groups_tables = `CREATE TABLE IF NOT EXISTS groups (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT,
			creator_id TEXT NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (creator_id) REFERENCES users(id)
		);`

	Group_members_tables = `CREATE TABLE IF NOT EXISTS groupmembers (
			id TEXT PRIMARY KEY,
			group_id TEXT NOT NULL,
			user_id TEXT NOT NULL,
			status TEXT CHECK( status IN ('pending', 'accepted') ) DEFAULT 'pending',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (group_id) REFERENCES groups(id),
			FOREIGN KEY (user_id) REFERENCES users(id)
		);`

	Notifications = `CREATE TABLE IF NOT EXISTS notifications (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,          
			content TEXT NOT NULL,             
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			read BOOLEAN DEFAULT FALSE,        
			FOREIGN KEY (user_id) REFERENCES users(id)
		);`

	Messages_table = `CREATE TABLE IF NOT EXISTS messages (
			id TEXT PRIMARY KEY,
			sender_id TEXT NOT NULL,
			recipient_id TEXT NOT NULL,
			content TEXT NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (sender_id) REFERENCES users(id),
			FOREIGN KEY (recipient_id) REFERENCES users(id)
		);`

	Posts_table = `CREATE TABLE IF NOT EXISTS posts (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			category TEXT NOT NULL,
			content TEXT NOT NULL,
			user_id TEXT NOT NULL,
			visibility TEXT CHECK( visibility IN ('public', 'private', 'limited') ) DEFAULT 'public',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			image_path TEXT,
			FOREIGN KEY (user_id) REFERENCES users(id)
		);`

	CommentLikes_table = `CREATE TABLE IF NOT EXISTS commentlikes (
			id TEXT PRIMARY KEY,
			comment_id TEXT NOT NULL,
			user_id TEXT NOT NULL,
			FOREIGN KEY (comment_id) REFERENCES comments(id),
			FOREIGN KEY (user_id) REFERENCES users(id)
		);`

	Comments_table = `CREATE TABLE IF NOT EXISTS comments (
			id TEXT PRIMARY KEY,
			post_id TEXT NOT NULL,
			content TEXT NOT NULL,
			user_id TEXT NOT NULL,
			username TEXT NOT NULL,
			created_at DATETIME NOT NULL,
			FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
		);`

	Unlikes_Table = `CREATE TABLE IF NOT EXISTS unlikes (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			post_id TEXT NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (post_id) REFERENCES posts(id)
		);`

	UnlikesComment_Table = `CREATE TABLE IF NOT EXISTS unlikescomment (
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		comment_id TEXT NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (comment_id) REFERENCES comment(id)
	);`
)
