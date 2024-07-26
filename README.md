<a id="readme-top"></a>

<!-- PROJECT LOGO -->
<br />
<div align="center">
<h1 align="center">PDF2Q</h1>

  <h3 align="center">
    An automation project to generate questions from a pdf using gpt4-omni
  </h3>
  <br/>
  <br/>
</div>

<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#contact">Contact</a></li>
    <li><a href="#acknowledgments">Acknowledgments</a></li>
  </ol>
</details>

<!-- ABOUT THE PROJECT -->

## About The Project

This project was made as an automation system for a teacher to generate questions from a pdf file using the openai api, the project uses tview as its text user interface (TUI), directory creation, json file creation, api requests, .env file, turning pdfs to images and turning them to base64 tosend over http and stuff. I hope it could be a life saver to someone who is starting go and doesent know how to use tview or send api requests.

### Built With

- [![Go][go.dev]][go-url]

<!-- GETTING STARTED -->

## Getting Started

First you have to have go installed on your system, I am developing on windows so my goto way is to use scoop.sh, visit the website and install scoop, then search for go and follow the instruction.

You also need poppler, I recommend installing that from scoop as well, we will need pdftocairo to turn pdfs to images.

### Prerequisites

Go has made dependencies so much easier, all you have to do after cloning the project is to run :

- Powershell:
  ```sh
  go mod tidy
  ```

The dependencies will be installed automagically.

### Installation

1. Get an api key from openai, it will cost a buck or five, and if you know of a free api, feel free to use it.

2. Clone the repo
   ```sh
   git clone https://github.com/mehybinator/pdf2q.git
   ```
3. Install go dependencies
   ```sh
   go mod tidy
   ```
4. Create a ".env" file and input the following :
   ```
   OPENAI_API = "YOUR_API_KEY";
   QUESTION_AMOUNT = "5";
   ```

<!-- USAGE EXAMPLES -->

## Usage

Run the program with :

```sh
make run
```

If you dont have make installed just use :

```sh
go build -o ./bin/pdf2q.exe ./cmd
./bin/pdf2q.exe
```

Running the code initially will cause an error, this is due to the fact that on initial run, the code will create the `data` folder with `images`, `pdfs` and `questions` subfolders, but the pdfs directory needs to contain a pdf for the program to work.

After you put your pdfs in the respective directory, running the program should show you a list of your pdfs, go ahead and pres enter on your desired pdf,
immediately the images should be generated and the api request should be sent, now you have to wait for gpt4-omni to respond, you should see both errors and success messages after the program is finished, now navigate to data/question and you should see a json file names rightly after your pdf file, this file contains the generated question.

<!-- CONTACT -->

## Contact

Mehran Arkak - mehran.arkak@protonmail.com

Project Link: [https://github.com/mehybinator/pdf2q](https://github.com/mehybinator/pdf2q)

<!-- ACKNOWLEDGMENTS -->

## Acknowledgments

- [Tview](https://github.com/rivo/tview)
- [Tcell](https://github.com/gdamore/tcell)
- [Poppler](https://poppler.freedesktop.org)
- [Golang](https://go.dev)

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->

[license-shield]: https://img.shields.io/github/license/mehybinator/pdf2q.svg?style=for-the-badge
[license-url]: https://github.com/mehybinator/pdf2q/blob/master/LICENSE.txt
[go.dev]: https://img.shields.io/badge/go-0769AD?style=for-the-badge&logo=go&logoColor=white
[go-url]: https://go.dev
