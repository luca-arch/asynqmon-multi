<!DOCTYPE html>

<html lang="en">
	<head>
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<title>{{ .Title }}</title>

		<style>
		:root {
			--bg: #ffffff;
			--text: #111111;
			--border: #cccccc;
		}

		@media (prefers-color-scheme: dark) {
			:root {
				--bg: #0f1115;
				--text: #e6e6e6;
				--border: #2a2f3a;
			}
		}

		html, body {
			background: var(--bg);
			height: 100%;
			margin: 0;
		}

		#menu {
			text-align: center;
		}

		#menu a {
			color: var(--text);
			font-family: monospace;
			font-size: 16px;
			font-weight: bold;
			padding: 1rem;
			text-decoration: none;
		}

		#menu a:hover {
			color: #6699CC;
			text-decoration: underline;
		}

		iframe {
			border: none 0;
			height: 100%;
			margin: 0;
			padding: 0;
			width: 100%;
		}
		</style>
	</head>

	<body>
		<div id="menu">
			{{range $name, $q := .Queues}}
				<a href="/{{ $name }}" target="asynqmon">{{ $name }} (db{{ $q.RedisDB }})</a>
			{{end}}
		</div>

		<iframe name="asynqmon" />
	</body>
</html>