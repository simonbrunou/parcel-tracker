type Translations = Record<string, string>;

const en: Translations = {
  // App
  "app.name": "Parcel Tracker",

  // Common
  "common.back": "Back",
  "common.cancel": "Cancel",
  "common.delete": "Delete",
  "common.save": "Save",
  "common.loading": "Loading...",
  "common.retry": "Retry",
  "common.edit": "Edit",
  "common.manual": "Manual",

  // Navbar
  "nav.toggleTheme": "Toggle theme",
  "nav.logout": "Logout",
  "nav.language": "Language",

  // Login
  "login.signIn": "Sign in to continue",
  "login.createPassword": "Create your password to get started",
  "login.password": "Password",
  "login.passwordPlaceholder": "Enter your password",
  "login.confirmPassword": "Confirm Password",
  "login.confirmPasswordPlaceholder": "Confirm your password",
  "login.signInButton": "Sign In",
  "login.getStarted": "Get Started",
  "login.passwordsMismatch": "Passwords do not match",
  "login.passwordTooShort": "Password must be at least 8 characters",
  "login.failed": "Login failed",

  // Dashboard
  "dashboard.title": "My Parcels",
  "dashboard.archive": "Archive",
  "dashboard.parcelCount": "{count} parcel",
  "dashboard.parcelCountPlural": "{count} parcels",
  "dashboard.addParcel": "Add Parcel",
  "dashboard.searchPlaceholder": "Search parcels...",
  "dashboard.allStatuses": "All statuses",
  "dashboard.noMatch": "No parcels match your filters",
  "dashboard.clearFilters": "Clear filters",
  "dashboard.noParcels": "No parcels yet",
  "dashboard.addFirst": "Add your first parcel to start tracking",

  // Add Parcel
  "addParcel.title": "Add Parcel",
  "addParcel.trackingNumber": "Tracking Number",
  "addParcel.trackingPlaceholder": "e.g. 1Z999AA10123456784",
  "addParcel.carrier": "Carrier",
  "addParcel.customName": "Custom Name",
  "addParcel.namePlaceholder": "e.g. New Headphones",
  "addParcel.notes": "Notes",
  "addParcel.notesPlaceholder": "Optional notes...",
  "addParcel.adding": "Adding...",
  "addParcel.submit": "Add Parcel",
  "addParcel.failed": "Failed to create parcel",

  // Parcel Detail
  "detail.refreshTracking": "Refresh tracking",
  "detail.archive": "Archive",
  "detail.unarchive": "Unarchive",
  "detail.tracking": "Tracking",
  "detail.carrier": "Carrier",
  "detail.status": "Status",
  "detail.estimatedDelivery": "Est. delivery",
  "detail.added": "Added",
  "detail.notes": "Notes",
  "detail.archived": "Archived",
  "detail.name": "Name",
  "detail.trackingNumber": "Tracking Number",
  "detail.trackingHistory": "Tracking History",
  "detail.addEvent": "+ Add Event",
  "detail.message": "Message *",
  "detail.messagePlaceholder": "e.g. Package arrived at sorting facility",
  "detail.eventStatus": "Status",
  "detail.location": "Location",
  "detail.locationPlaceholder": "Optional",
  "detail.submitEvent": "Add Event",
  "detail.deleteTitle": "Delete parcel?",
  "detail.deleteMessage": "This will permanently delete this parcel and all its tracking events.",
  "detail.loadFailed": "Failed to load parcel",

  // Timeline
  "timeline.empty": "No tracking events yet",

  // Statuses
  "status.unknown": "Unknown",
  "status.info_received": "Info Received",
  "status.in_transit": "In Transit",
  "status.out_for_delivery": "Out for Delivery",
  "status.delivered": "Delivered",
  "status.failed": "Failed",
  "status.expired": "Expired",

  // Dashboard errors
  "dashboard.loadFailed": "Failed to load parcels",

  // Not Found
  "notFound.title": "Page not found",
  "notFound.message": "The page you're looking for doesn't exist.",
  "notFound.goHome": "Go to Dashboard",

  // Archive confirmation
  "detail.archiveTitle": "Archive parcel?",
  "detail.archiveMessage": "This parcel will be moved to the archive.",
  "detail.unarchiveTitle": "Unarchive parcel?",
  "detail.unarchiveMessage": "This parcel will be restored from the archive.",

  // Toasts
  "toast.parcelCreated": "Parcel added successfully",
  "toast.parcelDeleted": "Parcel deleted",
  "toast.parcelArchived": "Parcel archived",
  "toast.parcelUnarchived": "Parcel unarchived",
  "toast.parcelUpdated": "Parcel updated",
  "toast.trackingRefreshed": "Tracking refreshed",
  "toast.eventAdded": "Event added",
  "toast.eventDeleted": "Event deleted",
  "toast.error": "An error occurred",

  // Relative time
  "time.justNow": "just now",
  "time.minutesAgo": "{n}m ago",
  "time.hoursAgo": "{n}h ago",
  "time.daysAgo": "{n}d ago",
};

const fr: Translations = {
  // App
  "app.name": "Suivi de Colis",

  // Common
  "common.back": "Retour",
  "common.cancel": "Annuler",
  "common.delete": "Supprimer",
  "common.save": "Enregistrer",
  "common.loading": "Chargement...",
  "common.retry": "R\u00e9essayer",
  "common.edit": "Modifier",
  "common.manual": "Manuel",

  // Navbar
  "nav.toggleTheme": "Changer de th\u00e8me",
  "nav.logout": "D\u00e9connexion",
  "nav.language": "Langue",

  // Login
  "login.signIn": "Connectez-vous pour continuer",
  "login.createPassword": "Cr\u00e9ez votre mot de passe pour commencer",
  "login.password": "Mot de passe",
  "login.passwordPlaceholder": "Entrez votre mot de passe",
  "login.confirmPassword": "Confirmer le mot de passe",
  "login.confirmPasswordPlaceholder": "Confirmez votre mot de passe",
  "login.signInButton": "Se connecter",
  "login.getStarted": "Commencer",
  "login.passwordsMismatch": "Les mots de passe ne correspondent pas",
  "login.passwordTooShort": "Le mot de passe doit contenir au moins 8 caract\u00e8res",
  "login.failed": "\u00c9chec de connexion",

  // Dashboard
  "dashboard.title": "Mes Colis",
  "dashboard.archive": "Archives",
  "dashboard.parcelCount": "{count} colis",
  "dashboard.parcelCountPlural": "{count} colis",
  "dashboard.addParcel": "Ajouter un Colis",
  "dashboard.searchPlaceholder": "Rechercher des colis...",
  "dashboard.allStatuses": "Tous les statuts",
  "dashboard.noMatch": "Aucun colis ne correspond \u00e0 vos filtres",
  "dashboard.clearFilters": "Effacer les filtres",
  "dashboard.noParcels": "Aucun colis",
  "dashboard.addFirst": "Ajoutez votre premier colis pour commencer le suivi",

  // Add Parcel
  "addParcel.title": "Ajouter un Colis",
  "addParcel.trackingNumber": "Num\u00e9ro de suivi",
  "addParcel.trackingPlaceholder": "ex. 1Z999AA10123456784",
  "addParcel.carrier": "Transporteur",
  "addParcel.customName": "Nom personnalis\u00e9",
  "addParcel.namePlaceholder": "ex. Nouveaux \u00e9couteurs",
  "addParcel.notes": "Notes",
  "addParcel.notesPlaceholder": "Notes facultatives...",
  "addParcel.adding": "Ajout...",
  "addParcel.submit": "Ajouter le Colis",
  "addParcel.failed": "\u00c9chec de cr\u00e9ation du colis",

  // Parcel Detail
  "detail.refreshTracking": "Actualiser le suivi",
  "detail.archive": "Archiver",
  "detail.unarchive": "D\u00e9sarchiver",
  "detail.tracking": "Suivi",
  "detail.carrier": "Transporteur",
  "detail.status": "Statut",
  "detail.estimatedDelivery": "Livraison estim\u00e9e",
  "detail.added": "Ajout\u00e9",
  "detail.notes": "Notes",
  "detail.archived": "Archiv\u00e9",
  "detail.name": "Nom",
  "detail.trackingNumber": "Num\u00e9ro de suivi",
  "detail.trackingHistory": "Historique de suivi",
  "detail.addEvent": "+ Ajouter un \u00e9v\u00e9nement",
  "detail.message": "Message *",
  "detail.messagePlaceholder": "ex. Colis arriv\u00e9 au centre de tri",
  "detail.eventStatus": "Statut",
  "detail.location": "Lieu",
  "detail.locationPlaceholder": "Facultatif",
  "detail.submitEvent": "Ajouter l'\u00e9v\u00e9nement",
  "detail.deleteTitle": "Supprimer le colis ?",
  "detail.deleteMessage": "Cela supprimera d\u00e9finitivement ce colis et tous ses \u00e9v\u00e9nements de suivi.",
  "detail.loadFailed": "\u00c9chec du chargement du colis",

  // Timeline
  "timeline.empty": "Aucun \u00e9v\u00e9nement de suivi",

  // Statuses
  "status.unknown": "Inconnu",
  "status.info_received": "Informations re\u00e7ues",
  "status.in_transit": "En transit",
  "status.out_for_delivery": "En livraison",
  "status.delivered": "Livr\u00e9",
  "status.failed": "\u00c9chou\u00e9",
  "status.expired": "Expir\u00e9",

  // Dashboard errors
  "dashboard.loadFailed": "\u00c9chec du chargement des colis",

  // Not Found
  "notFound.title": "Page introuvable",
  "notFound.message": "La page que vous recherchez n'existe pas.",
  "notFound.goHome": "Aller au tableau de bord",

  // Archive confirmation
  "detail.archiveTitle": "Archiver le colis ?",
  "detail.archiveMessage": "Ce colis sera d\u00e9plac\u00e9 dans les archives.",
  "detail.unarchiveTitle": "D\u00e9sarchiver le colis ?",
  "detail.unarchiveMessage": "Ce colis sera restaur\u00e9 depuis les archives.",

  // Toasts
  "toast.parcelCreated": "Colis ajout\u00e9 avec succ\u00e8s",
  "toast.parcelDeleted": "Colis supprim\u00e9",
  "toast.parcelArchived": "Colis archiv\u00e9",
  "toast.parcelUnarchived": "Colis d\u00e9sarchiv\u00e9",
  "toast.parcelUpdated": "Colis mis \u00e0 jour",
  "toast.trackingRefreshed": "Suivi actualis\u00e9",
  "toast.eventAdded": "\u00c9v\u00e9nement ajout\u00e9",
  "toast.eventDeleted": "\u00c9v\u00e9nement supprim\u00e9",
  "toast.error": "Une erreur est survenue",

  // Relative time
  "time.justNow": "\u00e0 l'instant",
  "time.minutesAgo": "il y a {n} min",
  "time.hoursAgo": "il y a {n} h",
  "time.daysAgo": "il y a {n} j",
};

const locales: Record<string, Translations> = { en, fr };
export const supportedLocales = [
  { code: "en", label: "English" },
  { code: "fr", label: "Fran\u00e7ais" },
];

function getInitialLocale(): string {
  const saved = localStorage.getItem("locale");
  if (saved && locales[saved]) return saved;

  const browserLang = navigator.language.split("-")[0];
  if (locales[browserLang]) return browserLang;

  return "en";
}

let locale = $state(getInitialLocale());

export function getLocale(): string {
  return locale;
}

export function setLocale(l: string): void {
  if (!locales[l]) return;
  locale = l;
  localStorage.setItem("locale", l);
  document.documentElement.setAttribute("lang", l);
}

export function t(key: string, params?: Record<string, string | number>): string {
  const translations = locales[locale] || locales["en"];
  let value = translations[key] || locales["en"][key] || key;
  if (params) {
    for (const [k, v] of Object.entries(params)) {
      value = value.replaceAll(`{${k}}`, String(v));
    }
  }
  return value;
}

export function getStatusLabel(status: string): string {
  return t(`status.${status}`);
}

export function getLocaleCode(): string {
  return locale;
}
