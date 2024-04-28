import { useEffect, useState } from 'react';
import Movie from '../components/Movie';
import { default as IMovie } from '../types/movie';
import MovieSearch from '../components/MovieSearch';
import MovieAdd from '../components/MovieAdd';
import getMovies from '../functions/getMovies';

export default function Movies() {
    const [error, setError] = useState('');
    const [movies, setMovies] = useState([] as IMovie[]);

    const movieAdded = (movie: IMovie) => {
        console.log("movies has been added", movie)
        if (movies.length > 10) {
            movies.pop()
        }

        setMovies([movie, ...movies])
    }

    useEffect(() => {
        (() => {
            getMovies('').then(setMovies).catch(setError)
        })();
    }, [])


    return (
        <>
            <div className="container">
                <div className="row">
                    <h1 className="text-center">My Movies</h1>
                </div>

                <div className="row mb-3 justify-around">
                    <div className="col-12 offset-0 col-sm-6">
                        <MovieAdd onAdd={movieAdded} onFail={setError}/>
                    </div>

                    <div className="col-12 mt-3 offset-0 mt-sm-0 col-sm-6">
                        <MovieSearch onFind={setMovies} onFail={setError} />
                    </div>
                </div>

                {error != '' &&
                    <div id="create-error">{error}</div>
                }

                {movies.length === 0 && 
                <div className="row">
                    <h4 className="col-12 text-center">
                        You have no saved movies.
                    </h4>
                </div>
                }
                <div className="row">
                    {movies.map((movie) => <Movie movie={movie} />)}
                </div>
            </div>
        </>
    )
}
