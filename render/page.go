package render

var page = `
<!DOCTYPE html>
<html>
	<head>
		<title>%s</title>
		<base target="_blank" />

		<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/github-markdown-css/5.2.0/github-markdown.min.css" integrity="sha512-Ya9H+OPj8NgcQk34nCrbehaA0atbzGdZCI2uCbqVRELgnlrh8vQ2INMnkadVMSniC54HChLIh5htabVuKJww8g==" crossorigin="anonymous" referrerpolicy="no-referrer" />
		<style>
			.markdown-body {
				max-width: 830px;
				margin: 10px 20px;
			}
			:root {
				--pre-background-color: rgb(246, 248, 250);
			}
			@media (prefers-color-scheme: dark) {
				:root {
					--pre-background-color: rgb(22, 27, 34);
				}
			}
			pre { background-color: var(--pre-background-color) !important; }

		</style>
		<script>
		  function __preventNavigation() {
				[...document.querySelectorAll('a')].forEach((a) => {
					a.onclick = ((evt) => {
						evt.preventDefault();
						window.__handleNavigation(a.href ?? '')
					})
				});
			}
		async function __update() {
			document.querySelector('body').innerHTML = await __getText()
			__preventNavigation()
		}
		window.onload = () => __preventNavigation()
		</script>
	</head>
	<body class="markdown-body">%s</body>
</html>
`
