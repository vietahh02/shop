$(".list-product").slick({
  infinite: true,
  slidesToShow: 4,
  slidesToScroll: 1,
  autoplay: true,
  arrows: true,
  prevArrow:
      "<button type='button' class='slick-prev pull-left'><i class='fa-solid fa-play fa-rotate-180' style='color: #ffffff;'></i></button>",
  nextArrow:
      "<button type='button' class='slick-next pull-right'><i class='fa-solid fa-play' style='color: #ffffff;'></i></button>",
  responsive: [
    {
      breakpoint: 860,
      settings: {
        slidesToShow: 3,
        slidesToScroll: 1,
      },
    },
    {
      breakpoint: 625,
      settings: {
        slidesToShow: 2,
        slidesToScroll: 1,
      },
    },
    {
      breakpoint: 480,
      settings: {
        slidesToShow: 1,
        slidesToScroll: 1,
      },
    },
  ],
});

$(".list-product-late").slick({
  infinite: true,
  slidesToShow: 4,
  slidesToScroll: 1,
  autoplay: true,
  arrows: true,
  prevArrow:
      "<button type='button' class='slick-prev pull-left'><i class='fa-solid fa-play fa-rotate-180' style='color: #ffffff;'></i></button>",
  nextArrow:
      "<button type='button' class='slick-next pull-right'><i class='fa-solid fa-play' style='color: #ffffff;'></i></button>",
  responsive: [
    {
      breakpoint: 860,
      settings: {
        slidesToShow: 3,
        slidesToScroll: 1,
      },
    },
    {
      breakpoint: 625,
      settings: {
        slidesToShow: 2,
        slidesToScroll: 1,
      },
    },
    {
      breakpoint: 480,
      settings: {
        slidesToShow: 1,
        slidesToScroll: 1,
      },
    },
  ],
});

$(".list-hot-deal").slick({
  infinite: true,
  slidesToShow: 1,
  slidesToScroll: 1,
  autoplay: true,
  arrows: true,
  prevArrow:
      "<button type='button' class='slick-prev pull-left'><i class='fa-solid fa-play fa-rotate-180' style='color: #ffffff;'></i></button>",
  nextArrow:
      "<button type='button' class='slick-next pull-right'><i class='fa-solid fa-play' style='color: #ffffff;'></i></button>",
});

$(".list-special-deal").slick({
  infinite: true,
  slidesToShow: 1,
  slidesToScroll: 1,
  autoplay: true,
  arrows: true,
  prevArrow:
      "<button type='button' class='slick-prev pull-left'><i class='fa-solid fa-play fa-rotate-180' style='color: #ffffff;'></i></button>",
  nextArrow:
      "<button type='button' class='slick-next pull-right'><i class='fa-solid fa-play' style='color: #ffffff;'></i></button>",
});

$(".list-product-feature").slick({
  infinite: true,
  slidesToShow: 3,
  slidesToScroll: 1,
  autoplay: true,
  arrows: true,
  prevArrow:
      "<button type='button' class='slick-prev pull-left'><i class='fa-solid fa-play fa-rotate-180' style='color: #ffffff;'></i></button>",
  nextArrow:
      "<button type='button' class='slick-next pull-right'><i class='fa-solid fa-play' style='color: #ffffff;'></i></button>",
  responsive: [
    {
      breakpoint: 860,
      settings: {
        slidesToShow: 2,
        slidesToScroll: 1,
      },
    },
    {
      breakpoint: 625,
      settings: {
        slidesToShow: 1,
        slidesToScroll: 1,
      },
    },
  ],
});

$(".list-banner").slick({
  infinite: true,
  slidesToShow: 1,
  slidesToScroll: 1,
  autoplay: false,
  arrows: false,
});
