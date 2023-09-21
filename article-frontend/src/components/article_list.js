// ArticlesList.js
import React, { useState, useEffect } from 'react';

const ArticlesList = ({articleAPI}) => {
  const [articles, setArticles] = useState([]);
  const [filters, setFilters] = useState({
    published: null,
    author_id: null,
    limit: 10,
    published_before: null,
    published_after: null
  });

  useEffect(() => {
    articleAPI.getArticles(filters).then(data => {
      setArticles(data);
    });
  }, [filters]);

  return (
    <div>
      {/* Possibly some filter UI here */}
      <ul>
        {articles.map(article => (
          <li key={article.id}>
            <div className='article-details'>
                <h3>Title: {article.title}</h3>
                <h5>Article ID: {article.id}</h5>
                <p>{article.content}</p>
                <button onClick={() => articleAPI.deleteArticle(article.id)}>Delete Article</button>
                <p>Is published: {article.published}</p>
            </div>
            {/* More details about the article, perhaps a link to a detailed view */}
          </li>
        ))}
      </ul>
    </div>
  );
};

export default ArticlesList;
