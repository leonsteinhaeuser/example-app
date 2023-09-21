import ArticlesAPI from './api/articles';
import CreateArticle from './components/form_create_article';
import ArticlesList from './components/article_list';
import ArticleDetails from './components/article_details';

const authors = ["bd957105-44a2-46f1-8e52-1c75779d6c88", "79f09b9f-86c6-4216-b7b6-044b3e86b512", "982d5bfd-1df8-4512-b233-c1e437063a25", "5acf6932-0c71-4794-9fdf-0dabe428b6c8"];

const articleAPI = new ArticlesAPI();

function App() {
  return (
    <div className="App">
      <CreateArticle authors={authors} articleAPI={articleAPI} />

      <p>-----------------------------------------------------</p>

      <ArticlesList articleAPI={articleAPI} />

      <p>-----------------------------------------------------</p>

      <ArticleDetails articleId="505646b3-4c9e-4ec4-9865-3f7805874b54" articleAPI={articleAPI} />
    </div>
  );
}

export default App;
