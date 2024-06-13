document.addEventListener("DOMContentLoaded", function() {


    const picturesContainer = document.getElementById("pictures_container");
    const controlNext = picturesContainer.querySelector(".control_next");
    const controlPrev = picturesContainer.querySelector(".control_prev");
    const pictures = picturesContainer.querySelectorAll("label");

    const url = new URL(window.location.href);
    const searchParam = url.searchParams.get('s');
    const comicIndexParam = url.searchParams.get('ci');

    document.getElementById("button_description").addEventListener("click", function (e) {
        e.preventDefault();
        var selectedComicID = document.querySelector('input[name="d"]:checked').value;
        window.location.href = url.origin + url.pathname + "?s=" + searchParam + "&d=" + selectedComicID + "&ci=" + currentIndex;
    });

    let currentIndex = 0

    if (comicIndexParam && comicIndexParam.length > 0) {
        const parsedIndex = parseInt(comicIndexParam);
        if (!isNaN(parsedIndex)) {
            currentIndex = parsedIndex;
        }
    }

    function showSlide(index) {
        pictures.forEach((picture, i) => {
            if (i === index) {
                picture.style.display = "block"
                picture.querySelector("input").checked = true
            } else {
                picture.style.display = "none"
            }
        });
    }

    function moveRight() {
        currentIndex = (currentIndex + 1) % pictures.length;
        showSlide(currentIndex);
    }

    function moveLeft() {
        currentIndex = (currentIndex - 1 + pictures.length) % pictures.length;
        showSlide(currentIndex);
    }

    controlNext.addEventListener("click", function(e) {
        e.preventDefault();
        moveRight();
    });

    controlPrev.addEventListener("click", function(e) {
        e.preventDefault();
        moveLeft();
    });

    showSlide(currentIndex);
});