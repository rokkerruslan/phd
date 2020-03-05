<template>
  <div class="container">
    <!-- START ARTICLE FEED -->
    <section class="articles">
      <div class="column is-8 is-offset-2">
        <!-- START ARTICLE -->
        <div class="card article">
          <div class="card-content">
            <div class="media">
              <div class="media-content has-text-centered">
                <p class="title article-title">Introducing a new feature for paid subscribers</p>
                <div class="tags has-addons level-item">
                  <span class="tag is-rounded is-info">@skeetskeet</span>
                  <span class="tag is-rounded">May 10, 2018</span>
                </div>
              </div>
            </div>
            <div class="content article-body">
              <p>Non arcu risus quis varius quam quisque. Dictum varius duis at consectetur lorem. Posuere sollicitudin
                aliquam ultrices sagittis orci a scelerisque purus semper. </p>
              <p>Metus aliquam eleifend mi in nulla posuere sollicitudin aliquam ultrices. In hac habitasse platea
                dictumst vestibulum rhoncus est pellentesque elit. Accumsan lacus vel facilisis volutpat. Non sodales
                neque sodales ut etiam.
                Est pellentesque elit ullamcorper dignissim cras tincidunt lobortis feugiat vivamus.</p>
              <h3 class="has-text-centered">How to properly center tags in bulma?</h3>
              <p> Proper centering of tags in bulma is done with class:
              <pre>level-item</pre>
              Voluptat ut farmacium tellus in metus vulputate. Feugiat in fermentum posuere urna nec. Pharetra convallis
              posuere morbi leo urna molestie.
              Accumsan lacus vel facilisis volutpat est velit egestas. Fermentum leo vel orci porta. Faucibus interdum
              posuere lorem ipsum.</p>
            </div>
          </div>
        </div>
        <!-- END ARTICLE -->
        <!-- START ARTICLE -->
        <div class="card article">
          <div class="card-content">
            <div class="media">
              <div class="media-center">
                <img src="http://www.radfaces.com/images/avatars/daria-morgendorffer.jpg" class="author-image"
                     alt="Placeholder image">
              </div>
              <div class="media-content has-text-centered">
                <p class="title article-title">Sapien eget mi proin sed 🔱</p>
                <p class="subtitle is-6 article-subtitle">
                  <a href="#">@daria</a> on February 17, 2018
                </p>
              </div>
            </div>
            <div class="content article-body">
              <p>
                Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et
                dolore magna aliqua. Accumsan lacus vel facilisis volutpat est velit egestas. Sapien eget mi proin sed.
                Sit amet mattis vulputate enim.
              </p>
              <p>
                Commodo ullamcorper a lacus vestibulum sed arcu. Fermentum leo vel orci porta non. Proin fermentum leo
                vel orci porta non pulvinar. Imperdiet proin fermentum leo vel. Tortor posuere ac ut consequat semper
                viverra. Vestibulum lectus mauris ultrices eros.
              </p>
              <h3 class="has-text-centered">Lectus vestibulum mattis ullamcorper velit sed ullamcorper morbi. Cras
                tincidunt lobortis feugiat vivamus.</h3>
              <p>
                In eu mi bibendum neque egestas congue quisque egestas diam. Enim nec dui nunc mattis enim ut tellus. Ut
                morbi tincidunt augue interdum velit euismod in. At in tellus integer feugiat scelerisque varius morbi
                enim nunc. Vitae suscipit tellus mauris a diam.
                Arcu non sodales neque sodales ut etiam sit amet.
              </p>
            </div>
          </div>
        </div>
        <!-- END ARTICLE -->
        <!-- START PROMO BLOCK -->
        <section class="hero is-info is-bold is-small promo-block">
          <div class="hero-body">
            <div class="container">
              <h1 class="title">
                <i class="fa fa-bell-o"></i> Nemo enim ipsam voluptatem quia.</h1>
              <span class="tag is-black is-medium is-rounded">
                                    Natus error sit voluptatem
                                </span>
            </div>
          </div>
        </section>
        <!-- END PROMO BLOCK -->
        <!-- START ARTICLE -->
        <div class="card article">
          <div class="card-content">
            <div class="media">
              <div class="media-center">
                <img src="http://www.radfaces.com/images/avatars/angela-chase.jpg" class="author-image"
                     alt="Placeholder image">
              </div>
              <div class="media-content has-text-centered">
                <p class="title article-title">Cras tincidunt lobortis feugiat vivamus.</p>
                <p class="subtitle is-6 article-subtitle">
                  <a href="#">@angela</a> on October 7, 2017
                </p>
              </div>
            </div>
            <div class="content article-body">
              <p>
                Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et
                dolore magna aliqua. Accumsan lacus vel facilisis volutpat est velit egestas. Sapien eget mi proin sed.
                Sit amet mattis vulputate enim.
              </p>
              <p>
                Commodo ullamcorper a lacus vestibulum sed arcu. Fermentum leo vel orci porta non. Proin fermentum leo
                vel orci porta non pulvinar. Imperdiet proin fermentum leo vel. Tortor posuere ac ut consequat semper
                viverra. Vestibulum lectus mauris ultrices eros.
              </p>
              <h3 class="has-text-centered">“Everyone should be able to do one card trick, tell two jokes, and recite
                three poems, in case they are ever trapped in an elevator.”</h3>
              <p>
                In eu mi bibendum neque egestas congue quisque egestas diam. Enim nec dui nunc mattis enim ut tellus. Ut
                morbi tincidunt augue interdum velit euismod in. At in tellus integer feugiat scelerisque varius morbi
                enim nunc. Vitae suscipit tellus mauris a diam.
                Arcu non sodales neque sodales ut etiam sit amet.
              </p>
            </div>
          </div>
        </div>
        <!-- END ARTICLE -->
      </div>

    </section>
  </div>
</template>

<script>
    export default {
        name: "Blog"
    }
</script>

<style scoped>
  .container {
    background: #F0F2F4;
  }

  .navbar.is-white {
    background: #F0F2F4;
  }

  .navbar-brand .brand-text {
    font-size: 1.11rem;
    font-weight: bold;
  }

  .hero-body {
    background-image: url(https://upload.wikimedia.org/wikipedia/commons/thumb/f/f6/Plum_trees_Kitano_Tenmangu.jpg/1200px-Plum_trees_Kitano_Tenmangu.jpg);
    background-position: center;
    background-size: cover;
    background-repeat: no-repeat;
    height: 500px;
  }

  .articles {
    margin: 5rem 0;
    /*margin-top: -200px;*/
  }

  .articles .content p {
    line-height: 1.9;
    margin: 15px 0;
  }

  .author-image {
    position: absolute;
    top: -30px;
    left: 50%;
    width: 60px;
    height: 60px;
    margin-left: -30px;
    border: 3px solid #ccc;
    border-radius: 50%;
  }

  .media-center {
    display: block;
    margin-bottom: 1rem;
  }

  .media-content {
    margin-top: 3rem;
  }

  .article, .promo-block {
    margin-top: 6rem;
  }

  div.column.is-8:first-child {
    padding-top: 0;
    margin-top: 0;
  }

  .article-title {
    font-size: 2rem;
    font-weight: lighter;
    line-height: 2;
  }

  .article-subtitle {
    color: #909AA0;
    margin-bottom: 3rem;
  }

  .article-body {
    line-height: 1.4;
    margin: 0 6rem;
  }

  .promo-block .container {
    margin: 1rem 5rem;
  }
</style>