# gowiki
A toy application written in Golang that allows you to create pages, store them in a folder on a "server", retrieve them, and/or edit them

It contains 2 packages:
* `io` which handles reading, writing and storing data(pages and their corresponding text)
* `server` which is a toy HTTP server that accepts some requests and returns some responses depending on what endpoints you go to access.

Available and supported endpoints are:
* `/view/:title` which responds with a view for the page that corresponds to the title in the slug
* `/new/` renders a very simple form that allows you create a page and save it
* `/save/` receives the data from the form and stores a page corresponding to the data in a `.txt` file
* `/edit/:title` renders a simple form pre-populated with the data from the page whose title corresponds to the slug

Happy tinkering! 🎉 
