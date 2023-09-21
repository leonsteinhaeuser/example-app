// ArticleDetails.js
import React, { useState, useEffect } from 'react';

const ArticleDetails = ({ articleId, articleAPI }) => {
  const [article, setArticle] = useState(null);

  useEffect(() => {
    articleAPI.getArticleById(articleId).then(data => {
      setArticle(data);
    });
  }, [articleId]);

  if (!article) return <div>Loading...</div>;

  return (
    <div>
      <div className='article-details'>
        <h3>Title: {article.title}</h3>
        <h5>Article ID: {article.id}</h5>
        <p>{article.content}</p>
        <button onClick={() => articleAPI.deleteArticle(article.id)}>Delete Article</button>
        <p>Is published: {article.published}</p>
      </div>
    </div>
  );
};

export default ArticleDetails;
