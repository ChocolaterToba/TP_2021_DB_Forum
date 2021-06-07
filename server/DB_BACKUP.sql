--
-- PostgreSQL database dump
--

-- Dumped from database version 12.5 (Ubuntu 12.5-0ubuntu0.20.04.1)
-- Dumped by pg_dump version 12.5 (Ubuntu 12.5-0ubuntu0.20.04.1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: case_insensitive; Type: COLLATION; Schema: public; Owner: postgres
--

CREATE COLLATION public.case_insensitive (provider = icu, deterministic = false, locale = 'und-u-ks-level2');


ALTER COLLATION public.case_insensitive OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: forums; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.forums (
    forumid integer NOT NULL,
    title character varying(100) NOT NULL,
    creator character varying(50) NOT NULL COLLATE public.case_insensitive,
    forumname character varying(50) NOT NULL,
    posts_count integer DEFAULT 0 NOT NULL,
    threads_count integer DEFAULT 0 NOT NULL
);


ALTER TABLE public.forums OWNER TO postgres;

--
-- Name: forums_forumid_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.forums_forumid_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.forums_forumid_seq OWNER TO postgres;

--
-- Name: forums_forumid_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.forums_forumid_seq OWNED BY public.forums.forumid;


--
-- Name: posts; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.posts (
    postid integer NOT NULL,
    parentid integer DEFAULT 0 NOT NULL,
    creator character varying(50) NOT NULL COLLATE public.case_insensitive,
    message text,
    isedited boolean DEFAULT false NOT NULL,
    threadid integer NOT NULL,
    created timestamp(0) without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.posts OWNER TO postgres;

--
-- Name: posts_postid_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.posts_postid_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.posts_postid_seq OWNER TO postgres;

--
-- Name: posts_postid_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.posts_postid_seq OWNED BY public.posts.postid;


--
-- Name: threads; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.threads (
    threadid integer NOT NULL,
    creator character varying(50) NOT NULL COLLATE public.case_insensitive,
    title character varying(100) NOT NULL,
    forumname character varying(50) NOT NULL,
    message text,
    created timestamp(0) without time zone DEFAULT now() NOT NULL,
    rating integer DEFAULT 0 NOT NULL,
    threadname character varying(50) NOT NULL
);


ALTER TABLE public.threads OWNER TO postgres;

--
-- Name: threads_threadid_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.threads_threadid_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.threads_threadid_seq OWNER TO postgres;

--
-- Name: threads_threadid_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.threads_threadid_seq OWNED BY public.threads.threadid;


--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    userid integer NOT NULL,
    username character varying(50) NOT NULL COLLATE public.case_insensitive,
    email character varying(50) NOT NULL COLLATE public.case_insensitive,
    fullname character varying(100),
    description text
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Name: users_userid_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.users_userid_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.users_userid_seq OWNER TO postgres;

--
-- Name: users_userid_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.users_userid_seq OWNED BY public.users.userid;


--
-- Name: votes; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.votes (
    voteid integer NOT NULL,
    username character varying(50) NOT NULL COLLATE public.case_insensitive,
    threadname character varying(50) NOT NULL,
    upvote boolean DEFAULT false NOT NULL
);


ALTER TABLE public.votes OWNER TO postgres;

--
-- Name: votes_voteid_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.votes_voteid_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.votes_voteid_seq OWNER TO postgres;

--
-- Name: votes_voteid_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.votes_voteid_seq OWNED BY public.votes.voteid;


--
-- Name: forums forumid; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.forums ALTER COLUMN forumid SET DEFAULT nextval('public.forums_forumid_seq'::regclass);


--
-- Name: posts postid; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.posts ALTER COLUMN postid SET DEFAULT nextval('public.posts_postid_seq'::regclass);


--
-- Name: threads threadid; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.threads ALTER COLUMN threadid SET DEFAULT nextval('public.threads_threadid_seq'::regclass);


--
-- Name: users userid; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users ALTER COLUMN userid SET DEFAULT nextval('public.users_userid_seq'::regclass);


--
-- Name: votes voteid; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.votes ALTER COLUMN voteid SET DEFAULT nextval('public.votes_voteid_seq'::regclass);


--
-- Data for Name: forums; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.forums (forumid, title, creator, forumname, posts_count, threads_count) FROM stdin;
9	Ruga sua earum itaque suo, oneri, nullam volui sectantur.	opus.MMEq1bhu715cJd	o6TzVBIEkxCCs	0	0
\.


--
-- Data for Name: posts; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.posts (postid, parentid, creator, message, isedited, threadid, created) FROM stdin;
\.


--
-- Data for Name: threads; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.threads (threadid, creator, title, forumname, message, created, rating, threadname) FROM stdin;
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (userid, username, email, fullname, description) FROM stdin;
155	scio.symYmh1ujuI5rv	es.Xam0ZmddrUCfP@suivis.org	Madison Thomas	Melius de tu dum requirunt fac. Ac lege o vis hos tu absurdissimum inplicaverant fulget. Delectari de modo his certissimus repente haec. Dicentium nutu. Loco ad augendo quam enervandam. Recessus cavis veni febris. Ergo. Ea disputando illum et, laudibus homines.
156	iubes.CG6AhZUUp1cijV	erit.5n6BZ6Vd7viIR@estac.org	Emma Robinson	Ex e gemitu drachmam significaret es suffragia a, dixi. Tum si tuo da at locutum. Neglecta cantus sunt ingesta. Praeire an gaudens. Colligo ne stet habiti. Sim me in.
157	didici.rRLYHm1Dj1F5j1	dum.pj9aH6vU7D5Ij@moveatdari.org	Emily Wilson	Modi ne facti per medice afficior christus. Cordi. Voce eum avide sese, aer, vox gaudiis his id. Has comitatur se nostri re sum. Tui est album dum.
159	o.ZD3aMmVDJ1ifJd	diei.z1kyzZ1171C5P@piamquot.org	James Miller	Infirma. Amemur interrogem excitentur re iniquitatibus audiam. Tuae meque nominatur mole. Lapsus amemur rei ac. Perturbatione. Ibi esto neque abs. Boni omnium audis delectationem vox tum. Transisse pugno dum.
160	ea.01306M1u7vicpU	ipsa.BV9AmZ1DJVFcP@facmulus.net	Natalie Thomas	Propter lucis approbet mei re sim. Sat si. Diu dixi ex curo tota de canto piae, imaginibus. Benedicere adpropinquet conperero certa inquiunt re novum efficeret illi.
162	ad.D99AzHVDruCf7U	esto.d9906zD1pU5ij@sparsistot.org	Abigail Thomas	Parvulus lucet vox mea tolerare dum cadunt. Diu. Campis sunt actiones maeroribus beatitudinis huius. Post conor eo. Diei. Placentes tot.
163	deo.gl9Y6H1uJUC5Ru	gaudens.n3KAZzv1Pdiip@metumveinfligi.com	Benjamin Wilson	Faciliter et saucium de escae agerem. Somnis colligenda quasi quodam ut pane has. Interstitio hac.
165	e.2c3YmMDVRU5C7u	utroque.2F3A661d7DFIp@istassuis.net	Charlotte Brown	Ubi avertat noscendum deo possunt maris at. Amo fuit aliquid has illos usui ex non imagines.
166	gusta.CX9AhzDvPUC571	simile.589yhz1dPDFcR@abpotuero.net	Emma Martinez	Tuam. Iubens pro. Omni demonstrata huius rei piae fieri subintrat pulchris fallaciam. Vita vi ob quae.
168	a.ColahZUDrdcIPu	heremo.5QlB6zvvRvfcp@positusagro.net	Ethan Smith	Ne altera sana. Ut iudex thesauro stupor, antequam, ac. Ago fit at vasis latet recognoscimus. Memor en eripietur ob ne desivero in. Bibendo. Des ei. Pleno eam. Sum dura petat ordinatorem scirem adlapsu oneri sudoris, nominis. Mira illud concurrunt abigo anima, vitaliter, et alium.
169	hic.7bKA6MuDj1cc7v	quia.pB9AM61UJ1CIj@eoquo.org	Zoey Thomas	Quid se sum si eo. Vere difficultatis sum vox nec, hos, spargant. Manes animam omni meo mei interpellante sinu. Haec talibus item ioseph oblivionis grex loco rem. O.
171	peragravi.943ah6v17uf5Jd	ab.k29Bh6dDJUfcp@suaese.net	Jayden Williams	In laudandum et me abiciam auditur, simul. Voce res vae. Amplius. An. Ea. Me ne. Eloquio vi. Lux. E noe.
172	nunc.44kyHm1171f57v	de.g4906Hu1p1Ifj@animamrationi.org	Michael Smith	Iustitiae non diu israel, tui. David scit tum suffragia alieni, pax. Viam ad quos accipiat res adducor nemo tu. Surgere. Eloquentiam se ubi separatum ibi rem, opus familiaritate ea. Campos quamvis os e filiorum. An dum maestitiae interpellante ne offeratur.
174	a.715YZH1uPdC5j1	tenebant.RvfAHhDD7Vfij@carohanc.net	Elijah Jones	Da et amo aqua laudis cum, contenti. Careo. Lugens hi noverunt augeret. Audit videbat pulvere dominos explorandi rapinam cubile.
175	vocasti.IVc0mZUUPu5ijv	vere.fDCBM6DUrv5cP@hominumusum.com	Sofia Williams	Praeter impium. Ita ut eo temptat. Id mel hi mallem, via hi.
177	febris.C6Ibmm1UPDIiJ1	munerum.fm5bh61dpUc5R@nedonasti.org	Ethan Garcia	Ideo fleo abs dare, cito. Noe vi res apparens. Moderationi ut fraternae oraturis vi. Excusatio erit quidem sinus eloquentiam. Ut officiis. Recordantes se homo re tamdiu teque insidiarum, mala ego. Num eum cum deum sensarum. Si grex ita toto fac corpore. Gloriatur multos ac a sensis antequam conferamus doctrinae, erat.
178	linguarum.nhiaM61DJ155JV	febris.2Z5y6zvUp1CI7@curre.net	Elizabeth Martin	Ait aditum ab agnoscerem ex exterminantes in cura adversitas. Dura officiis meo duobus magna quot ad gusta sub. Sed surgere salvus orare. E et nec sicuti iam in, fide pertendam. Sum nam cupiant insania o ista pugno si inlecebrosa. An id ita interroget odore. Meo os diceretur.
180	vi.Ef5A6Z11RDF57u	e.QF50h6V1Jd5Cr@neflete.com	Elizabeth Williams	Universus dissimilia hi velim rebellis quaerentes, in. Ruga seu vicinior os nostram ac mel. Qua.
181	sociorum.Hxi06611p1IFJU	superbis.mXcB6MvvPuFcR@voluiea.net	Emma Moore	An ego. Mansuefecisti ambiendum beatum ait hi, viva. Solam ad ad ad ad sane. Nos casu genuit dominum, recordando tuas. Iustificas incipio mei. Qui habites de vocasti eo dispulerim. E quod faciunt varia soli se corpore pedes formosa. Re est novi me.
183	a.IQ5AhZvuJvF5jV	iohannem.FocyHMUuR1F5P@potuitmulus.net	Ella Robinson	Genuit maneas numerans malitia an audire. Scio ex exultatione e sed. Animo quicquam me curo aliae eam contenti idem. Iam rogantem. Odor commendavi antiqua a iam ut tu nos. Canenti. Super sint contendunt. Humanae. Infirmitas sive conantes demonstratus et scis insaniam inhaeseram.
184	ut.7y5ymHDvrv5i7u	ventris.JAc0M6U1jDF5P@nostrummiris.org	Emily Garcia	Ne resolvisti gaudebit his os diei ex modis. Peccato suo retarder colligantur. Ab fiat loca re, es. Se eis caro.
186	det.H2I0z6117UFFRU	quo.625ahHUuRDcc7@lunamnoe.net	David Robinson	Pede ceterarumque vitae fallax eo, tua. Notatum religione vi passionis unus mea, quaestionum intellecta locus. Nonnullius quas capacitas tibi saepe. Lene.
187	lunam.anfB6mVdpD5Ijv	mandamus.04I06zvDR1Icj@seamplum.org	Aiden Jones	Quaesiveram in elapsum direxi. In dei transcendi abs dinoscens vis ac sanare. Deo eos potuere occultum amo. Respuimus rem custodiant dum periculorum. Qui sic et sub meruit te. Motus teneo litteratura ergo. Agro hoc des mortui animos inde a. Pacem pro sacrificatori haberet ac.
189	num.pu8YmM11rvIC7u	ac.pU8yZZ1dP1i5R@acfiant.com	Ethan Moore	Viva manifestari venatio. Facie te aliquid requiro movent, velut. Respice. Salute nutantibus munere os nostrae dissimile distorta sub. Filiis adducor. Loqueretur dixi.
190	eo.5VsaMhv1J1iIrV	o.fUtYMZuU71FC7@augebistum.org	Chloe Thomas	Vos fac. Quaerunt venit en intrinsecus, carthaginem, tot. Inventus. Abs. Caput eo si illo, modi. Sensusque gustavi os fructus, laudem mutaveris magnifico volo sicubi.
192	fui.Eh8AM6UDrVfF7d	has.WZsy66DUjdc5p@atex.com	Joseph Davis	Aspernatione mei aer quia eris tu tantarum hymnum veritatem. Sapientiorem. Credit aures.
193	a.19sa6mVujVc5pd	canem.ukX0zmUUj15fj@dieimystice.net	Sophia Jones	Ullis oculi dum en scribentur nunc subire. Paratus tuo os motus contristari vita te laudem. Omni petitur post. Fine omni re mala deinde. Iugo carentes assuescere hoc, eius tunc. Malo tu an. Et nati e es sola ea cui, faciet bene. Abs alis oraturis odorum facies tot coloratae omnem. Soli mei est beatos.
195	finis.dCSY6Zuu7ucfjv	voluntas.1cSb6HuvR1fCR@quamdiure.org	Matthew Thomas	Modus temptatio. Os temptatum id ac recti solitis.
196	os.8F80m6UDpdIFjU	secreta.XIsYMHvD7d5FJ@reoceani.org	Chloe Thomas	En iterum vicinior. Euge etiamne consonant sonare interest sanum ea das unum. Reperiret ratio facile spem medicus medicus fortius. Eo servo patienter se autem tutor. Audivimus iube nosti die. Ea consilium caro eram dum vos tenacius. Interfui cor.
198	compagem.xtTYmmvu715fRu	nesciam.txXbmHuu7vI57@malusgrex.org	Chloe Jones	Sensus palpa vae. E adprobandi pane. An peccavit. Manu mea iugo agit equus, ait. Esse ac suo quis re. O congruentem. Ac aequo valeant vana eum. Olet usum una o sat inferiore. Ei nam o ait nota quo at e.
199	reptilia.UW8ah6dv7UCiR1	pro.Uw8AhMuDPUiiP@paratuslicet.org	Mason Miller	Praeter. Ut. Ei ex odoratus castrorum accende. Hae veritatem tum me, aut. Fierem fui perditum.
201	imitanti.d0TyZhDUpVCF71	nec.1YSAz6UVJ1F5r@sicorones.org	Benjamin Robinson	Adipisci. Serie adparet a caro nec vocatur fit lucis magnifico. Supervacuanea soni gaudere erogo das munda vi. Es tam liber sed verbis quaestionum dixi subdita privatio. Nec velut mentiatur ullo ipsaque hos, moveri, vox ac. Quem ut sum ipsaque incertus diu.
202	captus.XA806HuvpdIfJu	casu.syTyH6uurVC5P@corei.net	Daniel Davis	Et te cupidatatium recessus at, tam. Re oleum videndo somnis manu ne vae. Lucet requiro genuit pax responderent. Casu filiis fuit benedicitur, ea mella oleum. Audio si fui ita me vanae. Die cor reperta. Alios tu elati ab vides sic. Huc at audeo a ab, mihi o usurpant os. Piam beares.
204	requies.o4xaMM11Ju5ij1	tenetur.OGSa6h11R15cp@viviteligam.com	James Jones	Vox exterminantes mare prout, ore. En regem florum augendo redimas ne vis. In ullo casto hac eo. Nemo fastu iniquitate amat. Quo adpellata quot.
205	iam.VRoA66UDRuCf7d	at.1RoyZhv1rU557@eades.com	Abigail Brown	Vix miser cum dicimus omnem.
207	cum.98UwuY6vrucijd	affectu.9SVqu0MdpdI57@usumadducor.org	Joshua Jackson	Ut item potu visco haec.
208	videndi.9kmW1AZ1jvciPd	nam.k3mq10m1jvFcj@inludien.com	Avery Robinson	Lectorem augebis. Oderunt. Hoc obliviscamur en gaudeant rogo nos antris ille. Futuri amorem. Ea doce et dimitti os, latet.
209	pertinet.4LhE1BM1PUCipU	digni.4KmwDamVjucIJ@etquo.net	Lily Anderson	Sicubi vi ut. Ita tua valerent magistro qui fluctus novi. Propositi ad. Praesto cogitur ponderibus fornax rogo aut, ibi, saepius. Molestias o. A illos vi innocentia fui caro ob. Quicquam sub amaremus.
211	gustandi.dTzqU0HV7D5cpu	salute.dsZQuBhV7d5fJ@quoquoconatus.com	Emily Wilson	Fructum olet sacramenti eis ruga nominamus num. Primatum fac. Fiducia hae odor unde non. Qui agnoscerent proferatur. Eo meminerunt miseria temptamur mala, tot nam sed. Iesus volui dare corrigendam nos, cor amor quaerebam non. Carthaginem da ob me minusve.
212	se.qtMOUahd7D5C7U	e.qshqva6VRDIfp@solatu.net	William Thomas	Lucerna iniuste affectus en agebam at fieri regem.
214	istae.OWzEVAmDrUF5RU	non.EwmQ1AH17v55p@istidei.com	Mia Jackson	Hoc. Meo seu hic appareat tua.
215	arguitur.ZYZeV06VpdFFpu	ita.6yHeDyZDpd5F7@fideegenus.org	Benjamin Davis	Direxi saepe nollem. Latis respondi. Es coniunctam nec. Tui. Hac lineas. Hic tua das.
217	potens.6gmod0zu71c5PV	ingredior.646ovBh1rVcI7@rogodeo.net	Charlotte Martin	Seu maerores a benedicis viva ad. Petat hos ideo tu, tria.
218	ipsas.22Mwu0HUrUF57v	scientiae.N2he1BzUpvFIR@gustaitidem.com	Abigail Miller	Unde aliquantulum quemadmodum o intrant ac. Suarum. Secum cedendo te divellit thesauri congruentem, ac tot. Amplum eis tua in.
220	deus.Jukqua6URDIcrD	vix.PULqub6UJUICJ@istamut.com	Joshua Martin	Sint coegerit. Praeteritorum antris ei ei, desidero. Capior filio unde. Sensum reconditum elian deerit abesset tui es viam ob.
221	aranea.819wVaHvJ1IFp1	ex.8U3ovbmujvicr@metumvedeerat.net	Madison Jackson	Mortuis eo et carnem inlexit quaeritur perdiderat. Cognoscere fructus amo ipsa, noe pulchras. Viderem at ob tuos sui. Ingesta satietate ac actione tacet aer et.
223	a.BZLe10HvP1c5PV	laudabunt.069eUYH1PdCFp@longetenendi.org	Joseph Thompson	Has utrumque facio tuorum, duabus sum, noe. Si subrepsit per totum rei ea, coram sibi. Istorum animales veris earum ergo mundum faciente officiis. Una de delectarentur modicum ait. Iam. Idem et. Da parte seu.
224	nitidos.kl3QUyzD7V55P1	fierem.9k9WuAHUjuCf7@fideivelut.org	Elijah Davis	Te oculo quietem nam, vitae libet. Obliviscamur. Propter audivimus scientiae parum mei urunt cor. Flete cognitor tu tu desidero retibus modi consumma, illuc.
226	ac.i8KQuaMUrvIIRv	superbis.5SLwVa6UP15ij@sihos.net	Andrew Martin	Ascendens ago an insania, rogeris instat proferrem si subiugaverant. Es mole solitis hac curo. Habitare ad. Adlapsu de hi quos deteriore. In fallitur repeterent ac violari audi gestat si et. Eo sumendi in omni cantu iustitiae audio, tuam. Dum rem unus has pronuntianti a. Sparsa sensus o.
227	e.JolEVB6V7Uff71	ad.jqKod061pUC5J@potestsiderum.net	Charlotte Anderson	Res durum pax. Diei cor negotium sim ab eo vel nec oceani. Singula dispersione necesse ingentibus nos eam comprehendet eliqua. Cum porto stupor poterimus imaginum amasti. Eas vel aliud. Aer via hi vidi vocant. Recti agnovi die adversus laetitia et. Das tui medius inmortalem, falsi. Incaute.
229	simplicem.iAKwd0zVJdI5jv	sese.F0lO1bz1J15iJ@nullamserie.com	Benjamin Harris	Alius nunc. Eum hoc. Es. Sim eos discurro subditus ei pergo.
230	at.729EV0MURvi571	abscondo.R4LEvyZ1PV5cP@deusfit.org	Emma Jackson	Amo bibo des. Tuum. E. Facies est istae tu vel delectarentur. Datur casu. Usum. Cogitetur quanto est probet, defrito se. Eo splendorem meo quaerit, re, ex et. Det laqueis mirum habet tenent ardentius num ea.
232	infinitum.1PiQDB61Rdfcj1	pede.1jCeda6dRuiF7@unainest.net	Isabella Jones	Nolo es quaestio in in ab quadam a o.
233	leve.TpiqVy6D7d5fpd	a.8PIq10h1j15fj@veroqua.net	Joshua Smith	Suavia alii e demetimur te an aurium es nuntiata. Confessione seu ipsas sed meminisse minusve dulce. Magnum o.
235	a.EUCQdAH17dCI7U	avertat.wd5eV0h1Jucfp@viduabus.org	Joshua Miller	Pulchritudine et. Quoquo occursantur inruentes sonos ego magnificet. Afficit os delectatio. Deliciosas tot reminiscentis res des. Modestis dum fierem ore noe tali grex adversitas proximi.
236	nollem.VhCo1b6VP1F5RV	vere.d6IW1aHvjDF5r@quantomelior.org	Liam Jackson	Cogitur manu in ambitiones socias afuerit. In spes illi alieno secum tui prout hae si. Hos ulla creator miles et meae. Sinis cedentibus ea amorem diutius. Cor illa maneas in intonas horrendum. Eos fallax tu da. Manducat amplius mentem. Affectent boni incideram si simus remisisti istas. Nova.
238	longe.MKIW106D71icJV	adest.m35wU06vJUF5P@eumdicimur.com	Elizabeth Jackson	Stat sono. Retibus o regem discernebat vel omnipotens, utrumque, serviam prodest. Sciri oceani ob hanc quisque.
239	eunt.OlIqd0ZuPV5ipv	rem.qLfQvymvPdfF7@venitsanabis.org	Emma Johnson	Ubi. Omnimodarum. Motus quo typho fine, abs. Multique augeret quae dico lenia passionum os tuus. Diligit grex quot respondes homo res loqueremur. Metum solebat. Vana muta o. Fiat constrictione o ei a illic regina sua, agit. Consulens habendum eo imperas agit.
241	o.n5CQuAMd7Df5Ru	opus.N5iod0ZDJd5cr@campissilente.com	Ethan Smith	Remisisti rursus qua bibendo metas ut. Foeda fit aranea. Animae hoc. Lata.
242	eis.9x5OdbZdPDFI7V	secreto.lSIW1Yzvr1Ifp@necantorare.org	Joshua Thomas	Candorem lux hi praeciderim os obliti fit. Hae te quoque mirabili dari eo, hi. Tum vim fastu meditatusque, ibi es. At re memoriae aut ei, placentes.
244	comitatur.9o5WVbzdp1CfPD	os.9qFwvbHU7ufIJ@nosda.com	Ava Miller	Sanguine paucis temporum fit, tu vox perpetret eum dominus. Pacem ideo sic animas mutant refrenare alibi, recolerentur. Piae te mediator lucis augeret aer. Vim sapit fias. An umbrarum vox audivi melior exserentes. Volvere deinde vi gyros. Cuius fac eum refero, sedet. Contremunt ne meminisse dicturus fames ibi sit nec.
245	admonente.2wCw1AM1pDi5Jd	est.NqfQD0zDJu55P@tanguntilli.org	Michael Jackson	Latet quanta. Tolerat ei quos tum. Aut. Aranea cadavere eo rem, an, se veritatem.
247	quanti.2a5oDB6U7d5Cj1	id.4BFoUaHuR1i57@cubiledoleat.net	Lily Taylor	Donec propinquius reddi avertat. Timeo via sit hi, veni, id. Quem et ioseph. Alta audi mali modis novum victor. Tum. Et cupientem proferrem. Sedet latine hi coniunctione tui. Erit vasis plenariam commemoro velle resisto cor manducandi da. Fugasti.
248	pristinum.oGFQDy6VJV5i7D	quantum.Q2iQ106uPVcij@detu.com	Michael Jones	Noe res fui exteriorum, ea. Re si melodias. Artibus vis vae ago ei socias filium vi. Duxi agatur mea an ambitione carne e mea cogitatione.
250	eos.AJSOU0z1Pu5cJ1	o.aJ8qd0mVP1C5p@tribuisturpis.com	Liam Harris	Consumma gressum intime sermo ei se. Alias occupantur ab ita hos. Decet at. Potentias nolo at cura in, gero hic. Esau sua considerabo ago da inlusionibus velle nescio. Manu memores hac saluti ex credendum. Det ut re. Ante liquide.
251	dicimus.huteuy617U55rU	alis.hD8wvAmVrD5cp@agitohuic.com	Elijah Jones	Maerere vivam contristor re corporales iudicante.
253	infirmior.I68EVAZDrDIfr1	tu.56sQvYmvPvccJ@oleatab.com	Michael Smith	De ab sed quaero boni auribus ad des.
254	an.4hSQVb6DJ1C5r1	partes.4m8WvAzuPuiiP@egostellas.net	Olivia Jackson	Et fac accepimus des an. Num strepitu os exciderat cur en da. Abditis fiat. Sensus pristinum ex veni, tamquam.
256	laudatur.A3xev0ZVJd5Fr1	meridies.0l8O1Bz17v557@solalanguor.org	Olivia Thompson	Recondo mali se quaeque altius, mira depromi cognoscendum. Tuae sua ex et, rem inconsummatus duxi, lege e. En hilarescit cor eorum medius pleno. Innocentia sat suo. Ulla sese neglecta delet creator e adduxerunt tuam omnibus. Una generis lucente. Des meae tetigi.
257	noe.Mf8w1YMUJ15cjD	canitur.mISE1YhVJvi5R@totaoblitum.net	Abigail Williams	Rapiunt an illud ad, hymnus. Consulunt ubique dicimur an tu interrumpunt laude responsa. Sapiat subire meo ea. Flatus potui mare fudi dominaris. A rapit cogitur hae. Ieiuniis mavult reminiscentis sobrios. Id utilitatem.
259	sua.M8XWDaHdJUfiPd	ait.6STW1ahUrVf5j@eiquem.org	Zoey Smith	Periculorum. Te futuri manet et. A nolit. Reprehensum vix sui faciam num at, discerno, et.
260	eas.a8sqUYZ1rVCcrD	teneor.asSQu0HuP155p@ullisrecolo.com	William Brown	Re claudicans ergo. Contrectatae ita noe vanae adpetere. Vindicavit memini re. Abundare re enim dulcedo certa.
262	patriam.QQtwD06DpuCcpd	o.oOXQDBmujuCcJ@dieieo.org	Chloe Johnson	Bona si mei credit ad mel redimas sparsis. Sono ea adesset tamquam deus, diverso. Ut. Offeretur unde inventa. Ac iam tuo viva quis fallit. Dei civium mulus dicens, ei captans sic ne, deo.
263	a.vA8ouymV7u557V	ambiendum.1AxqVAh17dFCp@tuamamo.org	Elizabeth Martin	Sic re.
265	e.VGseUy6U71CcRU	sonorum.1NTe1BmU7UFFJ@unusetsi.com	Liam Martin	Quoquo de habes. Dura sentiens sit removeri disputante pro utrum ita ille. Efficeret forma quantum homines, senectute. Ne coapta. Dicis huiuscemodi cogo qua.
266	e.T48EU0MdPv55r1	succurrat.t4XWVbz17ucfR@autporto.net	Daniel White	Pulchra longum vim mirificum credendum tunc pax se eloquentes. Te. Da dei pro quaerebant sedentem. Cavis da res plenus ac, alioquin immortalis. Persequi propinquius ex fugam. Tum voluptatem mei re utinam sum metumve lucente ac. Illi sui. Suae digni has spe una una da placuit. Solo amo.
267	iudex.PRWOuyhU71cFJL	eant.J7eeDbMur1f5Rm@tutunc.net	Chloe Johnson	Iterum. En vos. Testis cui modico illius, doce. Et laudatur. Imposuit esse cupimus vestibus.
268	res.BJEQVyHdp1IfpV	se.hVowDyMDp15fj@reipaucis.org	Mason Thomas	Item in multi fit meminerunt. Ac pars gaudiis aqua augeret. Da seu lux tenebrae, quendam.
269	opus.MMEq1bhu715cJd	ponere.6MWQu0ZDrui5p@quiaprimo.net	Noah Brown	Ex cantu imprimi factumque vestibus extingueris eunt odorum habes. Praeposita generalis auri consilium res sat. O amemur canenti. Cur. Interroges vix eas.
270	quaerit.w3eqvb61RUf57u	ut.QleW1yZU7VcfJ@citose.com	Маркиз О-де-Колóн	Бездельник третьего разряда 😋
271	ibi.5IEQUYhDpufFJV	urunt.iFowVAZVpUIIP@suntnam.net	Anthony Moore	Tria religione ei retenta. Tot eas apparet vi piam peritia odium da tali. Valeret tuo se potui, en permissum, sim in. Animalibus temptat gemitu inperturbata te fine refrenare isto a. Figmentis.
272	rogo.6TwQVYmur15CPU	humanae.M8WWVY6djd5Cr@istilucem.net	Chloe Jones	Pulsatori vis amo fierem ad.
\.


--
-- Data for Name: votes; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.votes (voteid, username, threadname, upvote) FROM stdin;
\.


--
-- Name: forums_forumid_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.forums_forumid_seq', 9, true);


--
-- Name: posts_postid_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.posts_postid_seq', 2, true);


--
-- Name: threads_threadid_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.threads_threadid_seq', 1, true);


--
-- Name: users_userid_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.users_userid_seq', 272, true);


--
-- Name: votes_voteid_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.votes_voteid_seq', 17, true);


--
-- Name: forums forums_pk_id; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.forums
    ADD CONSTRAINT forums_pk_id PRIMARY KEY (forumid);


--
-- Name: forums forums_un_forumname; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.forums
    ADD CONSTRAINT forums_un_forumname UNIQUE (forumname);


--
-- Name: posts posts_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.posts
    ADD CONSTRAINT posts_pk PRIMARY KEY (postid);


--
-- Name: threads threads_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.threads
    ADD CONSTRAINT threads_pk PRIMARY KEY (threadid);


--
-- Name: threads threads_un_threadname; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.threads
    ADD CONSTRAINT threads_un_threadname UNIQUE (threadname);


--
-- Name: users users_pk_id; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pk_id PRIMARY KEY (userid);


--
-- Name: users users_un_email; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_un_email UNIQUE (email);


--
-- Name: users users_un_username; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_un_username UNIQUE (username);


--
-- Name: votes votes_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.votes
    ADD CONSTRAINT votes_pk PRIMARY KEY (voteid);


--
-- Name: votes votes_un_user_thread; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.votes
    ADD CONSTRAINT votes_un_user_thread UNIQUE (username, threadname);


--
-- Name: forums_forumname_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX forums_forumname_idx ON public.forums USING btree (forumname);


--
-- Name: threads_creator_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX threads_creator_idx ON public.threads USING btree (creator);


--
-- Name: threads_forumname_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX threads_forumname_idx ON public.threads USING btree (forumname);


--
-- Name: forums forums_fk_username; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.forums
    ADD CONSTRAINT forums_fk_username FOREIGN KEY (creator) REFERENCES public.users(username) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: posts posts_fk_threadid; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.posts
    ADD CONSTRAINT posts_fk_threadid FOREIGN KEY (threadid) REFERENCES public.threads(threadid) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: posts posts_fk_username; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.posts
    ADD CONSTRAINT posts_fk_username FOREIGN KEY (creator) REFERENCES public.users(username) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: threads threads_fk_forumname; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.threads
    ADD CONSTRAINT threads_fk_forumname FOREIGN KEY (forumname) REFERENCES public.forums(forumname) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: threads threads_fk_username; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.threads
    ADD CONSTRAINT threads_fk_username FOREIGN KEY (creator) REFERENCES public.users(username) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: votes votes_fk_threadname; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.votes
    ADD CONSTRAINT votes_fk_threadname FOREIGN KEY (threadname) REFERENCES public.threads(threadname) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: votes votes_fk_username; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.votes
    ADD CONSTRAINT votes_fk_username FOREIGN KEY (username) REFERENCES public.users(username) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

