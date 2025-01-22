package main

import(
    "bufio"
    "io"
    "os"
    "fmt"
)

type element_t int;
type root_t int;


const (
    H1 element_t = iota
    P
);

type tag_level struct{
    tag element_t
    text string
    level int
}


func read_html_file(path string) io.Reader{

    var reader io.Reader;
    file, err := os.Open(path);
    if err != nil {
        fmt.Printf("HTML file does not exist\n");
        return nil;
    }

    defer file.Close();

    reader = bufio.NewReader(file);
    return reader;
}

func generate_html(body_elements []tag_level, header_elements []tag_level) (string, error){

    html := "<html>\n";

    if header_elements != nil{
        html += " <header>\n";
        header, err := generate_sub(header_elements);
        if err != nil{
            return "", err;
        }

        html += header + "\n </header>\n";
    }

    if body_elements != nil{
        html = " <body>\n";
        body, err1 := generate_sub(body_elements);
        if err1 != nil{
            return "", err1;
        }

        html += body + "\n </body>\n";
    }
    html += "</html>";

    return html, nil;
}

func generate_sub(elements []tag_level) (string, error){

    var html string;
    var tag string;
    for _, element := range elements{
        spacing :=" "
        indent := element.level;
        for indent > 0 {
            spacing +=" ";
            indent -= 1;
        }


        switch(element.tag){
            case H1:
                tag = "H1";

            case P:
                tag = "P";
            
            default:
                return html, fmt.Errorf("Mismatch while generating html\n");
        }

        html += spacing + "<" + tag + ">\n" + element.text +"\n" + spacing + "</" + tag + ">";
    }

    return html, nil
}
        
