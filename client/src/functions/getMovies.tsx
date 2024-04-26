import { default as IMovie } from '../types/movie'

export default function getMovies(nameSearch: string): Promise<IMovie[]> {
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