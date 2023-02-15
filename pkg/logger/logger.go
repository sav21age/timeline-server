package logger
 
import (
   "github.com/rs/zerolog"
   "os"
)

var Logger zerolog.Logger

func InitLog()  {
   consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
   multiWriter := zerolog.MultiLevelWriter(consoleWriter, os.Stdout)
   Logger = zerolog.New(multiWriter).With().Timestamp().Logger()
}


// zerolog.SetGlobalLevel(zerolog.DebugLevel)
// zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
// log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339Nano})


// https://replit.com/talk/share/Zerolog-Golang-Example/80559
// package main

// import (
// 	"errors"
// 	"flag"
// 	"fmt"
// 	"os"

// 	"github.com/rs/zerolog"
// 	"github.com/rs/zerolog/log"
// )

// // Estructura de ejemplo
// type person struct {
// 	name string
// 	age  int
// }

// func main() {
// 	// zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
// 	debug := flag.Bool("debug", false, "sets log level to debug")

// 	flag.Parse()

// 	// Run sin el flag -debug
// 	zerolog.SetGlobalLevel(zerolog.InfoLevel)
// 	if *debug {
// 		// Run con el flag -debug
// 		zerolog.SetGlobalLevel(zerolog.DebugLevel)
// 	}

// 	// Para el manejo de errores no contemplados
// 	defer func() {
// 		if err := recover(); err != nil {
// 			strErr := fmt.Errorf("%v", err)

// 			log.Fatal().
// 				Err(strErr).
// 				Msgf("Error %s", "not contemplated")
// 		}
// 	}()

// 	// Set estructura de ejemplo
// 	person := person{name: "Joe", age: 84}
// 	// Log estructura de ejemplo (no es recomendable usar sprint, usar MarshalZerologObject)
// 	strPerson := fmt.Sprintf("%#v", person)
// 	log.Debug().
// 		Str("person", strPerson).
// 		Msg("This message appears only when log level set to Debug")

// 	log.Info().Msg("This message appears when log level set to Debug or Info")

// 	if e := log.Debug(); e.Enabled() {
// 		// Se realizan calculos y demas solo cuando debug esta activado.
// 		value := person.name
// 		e.Str("Name", value).Msg("some debug message")
// 	}

// 	// Error al leer un archivo
// 	file, errFile := os.Open("file+go")
// 	data := make([]byte, 100)
// 	count, _ := file.Read(data)
// 	if errFile != nil {
// 		log.Error().Err(errFile).Msgf("")
// 	} else {
// 		fmt.Printf("read %d bytes: %q\n", count, data[:count])
// 	}

// 	// Error contemplado lanzado a modo informativo
// 	err := errors.New("A contemplated error")
// 	service := "myservice"
// 	log.Warn().
// 		Err(err).
// 		Str("service", service).
// 		Msgf("Cannot start %s", service)

// 	// Error no contemplado al dividir por cero
// 	valZero := 0
// 	result := 10 / valZero
// 	fmt.Println(result)
// }
