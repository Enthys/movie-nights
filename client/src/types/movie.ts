export default interface Movie {
  id: number;
  name: string;
  description: string;
  imdbLink: string;
  genres: string[];
  avatar: string;
  imdbRating: number;
}