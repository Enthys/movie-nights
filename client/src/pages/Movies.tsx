import { useState } from 'react';
import Movie from '../components/Movie';
import { default as IMovie } from '../types/movie';

function getMovies(nameSearch: string): Promise<IMovie[]> {
    return new Promise<IMovie[]>((resolve, reject) => {
        fetch("/api/movies" + (nameSearch !== '' ? `?name=${nameSearch}` : ''))
            .then((resp) => {
                if (resp.status !== 200) {
                    reject(new Error("failed to retrieve movies"))
                }

                return resp.json();
            })
            .then((data) => {
                resolve(data.movies);
            })
    })
}

export default function Movies() {
    const [error, setError] = useState('');
    const [isSet, setIsSet] = useState(false)
    const [movies, setMovies] = useState([] as IMovie[]);
    const [search, setSearch] = useState('')

    if (!isSet) {
        getMovies('').then((movies) => {
            setMovies(movies);
            setIsSet(true);
        }).catch(setError);
    }

    return (
        <>
            <div className="container">
                <div className="row">
                    <h1 className="text-center">My Movies</h1>
                </div>

                <div className="row mb-3 justify-around">
                    <div className="col-12 offset-0 col-sm-6">
                        <form className="input-group" hx-ext="response-targets">
                            <input
                                name="name"
                                type="text"
                                className="form-control"
                                placeholder="Search for movie"
                                value={search}
                                onChange={(event) => setSearch(event.target.value)}
                            />
                            <button
                                type="submit"
                                className="btn btn-outline-success"
                                onClick={() => getMovies(search).then(setMovies)}
                            >Search</button>
                        </form>
                    </div>

                    <div className="col-12 mt-3 offset-0 mt-sm-0 col-sm-6">
                        <form className="input-group" hx-ext="response-targets">
                            <input name="url" type="text" className="form-control" placeholder="Movie IMDb link" />
                            <button type="submit" className="btn btn-outline-success">Add</button>
                        </form>
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
