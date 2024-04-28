import IMovie from '../types/movie';

export default function addMovie(link: string): Promise<IMovie> {
    return new Promise((resolve, reject) => {
        fetch('/api/movies', {
            method: 'POST',
            body: JSON.stringify({ link }),
            headers: {
                'Content-Type': 'application/json',
            },
        })
        .then(resp => {
            resp.json().then((data: { movie: IMovie, errors?: { link: string} }) => {
                console.log(data)
                if (data.errors !== undefined) {
                    reject(data.errors.link)
                    return
                }

                resolve(data.movie)
            })
        })
    })
}