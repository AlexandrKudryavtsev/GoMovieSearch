import { launch } from 'puppeteer';
import fs from 'fs/promises';
import path from 'path';

export const parseArgs = () => {
    const args = process.argv.slice(2);
    const destinationIndex = args.indexOf('--destination');

    if (destinationIndex === -1) {
        console.error('Usage: node app.js --destination <path to destination>');
        process.exit(1);
    }

    return {
        destination: args[destinationIndex + 1],
    };
}

const handleCaptchaRedirect = async (page) => {
    await page.waitForNavigation({
        waitUntil: 'networkidle2',
        timeout: 5000
    }).catch(() => { });

    const currentUrl = page.url();
    if (currentUrl.includes('/showcaptcha')) {
        console.log('Captcha detected! Please solve it manually...');
        await page.waitForNavigation({
            waitUntil: 'networkidle2',
            timeout: 0
        });
        console.log('Captcha solved, continuing parsing...');
    }
};

const parseMovies = async (kinopoiskListURL, totalPages = 1) => {
    const browser = await launch({ headless: false });
    const page = await browser.newPage();

    await page.goto("https://www.kinopoisk.ru");
    await handleCaptchaRedirect(page);

    const moviesData = [];

    for (let currentPage = 1; currentPage <= totalPages; currentPage++) {

        console.log(`Parsing page ${currentPage} of ${totalPages}...`);

        let baseUrl = kinopoiskListURL.split('?')[0];
        let pageUrl = currentPage === 1
            ? baseUrl
            : `${baseUrl}?page=${currentPage}`;

        try {
            await page.goto(pageUrl, {
                waitUntil: 'networkidle2',
                timeout: 30000
            });

            await handleCaptchaRedirect(page);

            const movieElements = await page.$$('div[data-test-id="movie-list-item"]');
            if (movieElements.length === 0) {
                console.log(`No movies found on page ${currentPage}, stopping parsing.`);
                break;
            }

            for (const movieEl of movieElements) {
                try {
                    const russianTitle = await movieEl.$eval(
                        '.styles_mainTitle__IFQyZ',
                        el => el.textContent.trim()
                    );

                    const originalTitle = await movieEl.$eval(
                        '.desktop-list-main-info_secondaryTitle__ighTt',
                        el => el.textContent.trim()
                    ).catch(() => null);

                    const yearText = await movieEl.$eval(
                        '.desktop-list-main-info_secondaryText__M_aus',
                        el => el.textContent.trim()
                    );
                    const year = yearText.match(/\d{4}/)?.[0] || null;

                    const detailsUrl = await movieEl.$eval(
                        'a[href^="/film/"]',
                        el => el.getAttribute('href')
                    );
                    const fullDetailsUrl = `https://www.kinopoisk.ru${detailsUrl}`;

                    const posterUrl = await movieEl.$eval(
                        '.styles_image__gRXvn',
                        el => el.getAttribute('src')
                    );
                    const fullPosterUrl = posterUrl.startsWith('//')
                        ? `https:${posterUrl}`
                        : posterUrl;

                    moviesData.push({
                        russianTitle,
                        originalTitle,
                        year,
                        detailsUrl: fullDetailsUrl,
                        posterUrl: fullPosterUrl,
                    });

                } catch (error) {
                    console.error(`Error parsing element on page ${currentPage}:`, error);
                }
            }
        } catch (error) {
            console.error(`Error loading page ${currentPage}:`, error);
            break;
        }
    }

    await browser.close();
    return moviesData;
};

const saveDataToJson = async (data, folderPath, fileName = 'data.json') => {
    try {
        const filePath = path.join(folderPath, fileName);

        await fs.mkdir(folderPath, { recursive: true });

        await fs.writeFile(filePath, JSON.stringify(data, null, 2));

        console.log(`Data successfully saved to ${filePath}`);
    } catch (err) {
        console.error('Error saving file:', err);
    }
}

(async () => {
    const { destination } = parseArgs();

    const KINOPOISK_LIST_URL = "https://www.kinopoisk.ru/lists/movies/top500/"
    const KINOPOISK_LIST_PAGES = 10

    try {
        const res = await parseMovies(KINOPOISK_LIST_URL, KINOPOISK_LIST_PAGES)

        saveDataToJson(res, destination)
    } catch (error) {
        console.error(error);
        await browser.close();
        process.exit(1);
    }
})();