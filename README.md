## Flags
- `-dimen :: <uint>`:

   	Avatar asset will be retrieved for dimension `0x01 << this` (default 10).

- `-opath :: <filepath>`:

   	Output path. If set, profile image will be retrieved and output to this file. Extensions will be appended for the appropriate MIME type.

- `-token :: <token>`

   	Bot token for Discord API.
- `-uid :: <id>`

  	User ID for which to fetch avatar.

- `-uname :: <name>[#<discriminator]`
		
		When this is given, the client starts an interactive mode prompting the user to log in , allowing the bot to introspect on their contacts list (`Relationship`s), performing a search for a username containing the given string and a matching discriminator where provided.
