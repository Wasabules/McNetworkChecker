import { writable, derived } from 'svelte/store';

export const locale = writable('fr');

export const LOCALES = [
  { code: 'fr', flag: '\u{1F1EB}\u{1F1F7}', label: 'FR' },
  { code: 'en', flag: '\u{1F1EC}\u{1F1E7}', label: 'EN' },
  { code: 'es', flag: '\u{1F1EA}\u{1F1F8}', label: 'ES' },
  { code: 'de', flag: '\u{1F1E9}\u{1F1EA}', label: 'DE' },
  { code: 'pt', flag: '\u{1F1F5}\u{1F1F9}', label: 'PT' },
];

// idx: 0=fr, 1=en, 2=es, 3=de, 4=pt
const S = {
  // App
  'app.title':        ['McNetworkChecker','McNetworkChecker','McNetworkChecker','McNetworkChecker','McNetworkChecker'],
  'app.subtitle':     ['Diagnostic Reseau Minecraft','Minecraft Network Diagnostic','Diagnostico de Red Minecraft','Minecraft Netzwerk-Diagnose','Diagnostico de Rede Minecraft'],

  // Steps
  'step.dns':         ['Resolution DNS','DNS Resolution','Resolucion DNS','DNS-Auflosung','Resolucao DNS'],
  'step.ping':        ['Ping ICMP','ICMP Ping','Ping ICMP','ICMP-Ping','Ping ICMP'],
  'step.traceroute':  ['Traceroute ICMP','ICMP Traceroute','Traceroute ICMP','ICMP-Traceroute','Traceroute ICMP'],
  'step.tcp':         ['Connexion TCP','TCP Connection','Conexion TCP','TCP-Verbindung','Conexao TCP'],
  'step.tcpTraceroute':['Traceroute TCP','TCP Traceroute','Traceroute TCP','TCP-Traceroute','Traceroute TCP'],
  'step.minecraft':   ['Serveur Minecraft','Minecraft Server','Servidor Minecraft','Minecraft-Server','Servidor Minecraft'],

  // Badges
  'badge.pending':    ['...','...','...','...','...'],
  'badge.running':    ['>>>','>>>','>>>','>>>','>>>'],
  'badge.ok':         ['OK','OK','OK','OK','OK'],
  'badge.fail':       ['FAIL','FAIL','FALLO','FEHLER','FALHA'],
  'badge.skip':       ['SKIP','SKIP','SKIP','SKIP','SKIP'],
  'badge.stopped':    ['---','---','---','---','---'],

  // Summaries
  'sum.avg':          ['moy','avg','med','Avg','med'],
  'sum.packets':      ['paquets','packets','paquetes','Pakete','pacotes'],
  'sum.hops':         ['sauts','hops','saltos','Hops','saltos'],
  'sum.players':      ['joueurs','players','jugadores','Spieler','jogadores'],

  // Buttons
  'btn.start':        ['Lancer','Start','Iniciar','Starten','Iniciar'],
  'btn.stop':         ['Arreter','Stop','Detener','Stoppen','Parar'],
  'btn.skip':         ['Skip >>','Skip >>','Skip >>','Skip >>','Skip >>'],
  'btn.copy':         ['Copier','Copy','Copiar','Kopieren','Copiar'],
  'btn.copied':       ['OK','OK','OK','OK','OK'],

  // Input
  'input.placeholder':['play.serveur.com ou 192.168.1.1:25565','play.server.com or 192.168.1.1:25565','play.servidor.com o 192.168.1.1:25565','play.server.com oder 192.168.1.1:25565','play.servidor.com ou 192.168.1.1:25565'],

  // Validation
  'val.portInvalid':  ['Port invalide ({port}) — doit etre entre 1 et 65535','Invalid port ({port}) — must be between 1 and 65535','Puerto invalido ({port}) — debe estar entre 1 y 65535','Ungultiger Port ({port}) — muss zwischen 1 und 65535 liegen','Porta invalida ({port}) — deve estar entre 1 e 65535'],
  'val.octet':        ['Adresse IP invalide — octet > 255','Invalid IP address — octet > 255','Direccion IP invalida — octeto > 255','Ungultige IP-Adresse — Oktett > 255','Endereco IP invalido — octeto > 255'],
  'val.tooLong':      ['Nom de domaine trop long','Domain name too long','Nombre de dominio muy largo','Domainname zu lang','Nome de dominio muito longo'],
  'val.spaces':       ['Les espaces ne sont pas autorises','Spaces are not allowed','Los espacios no estan permitidos','Leerzeichen sind nicht erlaubt','Espacos nao sao permitidos'],
  'val.badChars':     ['Caracteres invalides','Invalid characters','Caracteres invalidos','Ungultige Zeichen','Caracteres invalidos'],
  'val.incomplete':   ['Nom de domaine incomplet (ex: play.serveur.com)','Incomplete domain (e.g. play.server.com)','Dominio incompleto (ej: play.servidor.com)','Unvollstandiger Domain (z.B. play.server.com)','Dominio incompleto (ex: play.servidor.com)'],
  'val.invalid':      ['Adresse invalide','Invalid address','Direccion invalida','Ungultige Adresse','Endereco invalido'],

  // Skip messages
  'skip.auto':        ['Etape non necessaire (adresse IP directe)','Step not needed (direct IP address)','Paso no necesario (direccion IP directa)','Schritt nicht notig (direkte IP-Adresse)','Etapa nao necessaria (endereco IP direto)'],
  'skip.user':        ['Etape passee par l\'utilisateur','Step skipped by user','Paso omitido por el usuario','Schritt vom Benutzer ubersprungen','Etapa ignorada pelo utilizador'],

  // Minecraft info
  'mc.protocol':      ['Protocole','Protocol','Protocolo','Protokoll','Protocolo'],

  // Report
  'report.title':     ['Rapport de diagnostic','Diagnostic Report','Informe de diagnostico','Diagnosebericht','Relatorio de diagnostico'],
  'report.copied':    ['Copie !','Copied!','Copiado!','Kopiert!','Copiado!'],

  // Paste / Upload
  'paste.upload':     ['Heberger','Upload','Subir','Hochladen','Enviar'],
  'paste.uploading':  ['Envoi...','Uploading...','Subiendo...','Hochladen...','Enviando...'],
  'paste.success':    ['Lien copie !','Link copied!','Enlace copiado!','Link kopiert!','Link copiado!'],
  'paste.error':      ['Erreur : {err}','Error: {err}','Error: {err}','Fehler: {err}','Erro: {err}'],
  'paste.open':       ['Ouvrir','Open','Abrir','Offnen','Abrir'],
};

const IDX = { fr: 0, en: 1, es: 2, de: 3, pt: 4 };

export const t = derived(locale, ($locale) => {
  const idx = IDX[$locale] ?? 1;
  return (key, params) => {
    const arr = S[key];
    let str = arr ? (arr[idx] || arr[1] || key) : key;
    if (params) {
      for (const [k, v] of Object.entries(params)) {
        str = str.replaceAll(`{${k}}`, v);
      }
    }
    return str;
  };
});
