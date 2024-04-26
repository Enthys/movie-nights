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
            .then(resp => resp.json())
            .catch(reject)
                .then((data: { movie: IMovie }) => resolve(data.movie))
                .catch(reject)
    })
}