package template

import "fmt"

func Index(url, title, cname string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta http-equiv="refresh" content="0; url=%[1]s" />
  <title>%[2]s</title>
</head>
<body>
  <p><a href="%[1]s">Click here</a> if not redirected automatically.</p>
</body>
</html>
`,
		url,
		title,
	)
}

func CNAME(url, title, cname string) string {
	return cname
}

func README(url, title, cname string) string {
	return fmt.Sprintf(`### %[2]s

Redirects to [%[1]s](%[1]s).

###### Automatically generated by [kshvmdn/redir](https://github.com/kshvmdn/redir).
`,
		url,
		title,
	)
}