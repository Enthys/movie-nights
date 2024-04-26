import { useState } from "react"
import { default as IMovie } from '../types/movie'
import getMovies from "../functions/getMovies"

interface MovieSearchArgs {
    onFind: (movies: IMovie[]) => void
    onFail: (err: string) => void
}

export default function MovieSearch({ onFind, onFail }: MovieSearchArgs) {
    const [search, setSearch] = useState('')

    return (
        <div className="input-group" hx-ext="response-targets">
            <input
                name="name"
                type="text"
                className="form-control"
                placeholder="Search for movie"
                value={search}
                onChange={(event) => setSearch(event.target.value)}
            />
            <button
                type="button"
                className="btn btn-outline-success"
                onClick={() => getMovies(search).then(onFind).catch(onFail)}
            >Search</button>
        </div>
    )
}