// Default images configuration
export const DEFAULT_IMAGES = {
    // Gallery categories
    gallery: {
        wedding: '/storage/weeding/photo_2026-01-10%2001.28.47.jpeg',
        marryme: '/storage/marryme/photo_2026-01-10%2001.21.31.jpeg',
        birthday: '/storage/birthday/photo_2026-01-10%2001.17.18.jpeg',
        baby_shower: '/storage/baby%20shower/photo_2026-01-10%2001.34.18.jpeg',
        bapteme: '/storage/bapteme/photo_2026-01-10%2001.31.37.jpeg',
        loveroom: '/storage/loveroom/photo_2026-01-10%2001.33.08.jpeg',
        congrats: '/storage/congrats/photo_2026-01-10%2001.25.09.jpeg'
    },
    // Rental/Location categories
    rental: {
        default: '/storage/location/Arche%20semi%20waves.jpeg',
        centerpiece: '/storage/location/Centre%20de%20table%20mariage.jpeg',
        backdrop: '/storage/location/Support%20toile.jpeg',
        flower: '/storage/location/Fleur%20blanche%20champagne.jpeg',
        other: '/storage/location/Chaise%20chiavari.jpeg'
    }
}

// Helper function to get default image
export function getDefaultImage(type, category) {
    if (type === 'gallery') {
        return DEFAULT_IMAGES.gallery[category] || DEFAULT_IMAGES.gallery.wedding
    } else if (type === 'rental') {
        return DEFAULT_IMAGES.rental[category] || DEFAULT_IMAGES.rental.default
    }
    return ''
}

// Helper function to get image URL with fallback
export function getImageWithFallback(imageUrl, type, category) {
    if (imageUrl && imageUrl.trim() !== '') {
        return imageUrl
    }
    return getDefaultImage(type, category)
}
